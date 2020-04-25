package pluginapp

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"time"

	"github.com/ironzhang/x-pearls/log"

	"git.ablecloud.cn/ablecloud/ac-comm-lib/pluginapp/configurator/jsonc"
)

var exit = os.Exit

// G 应用框架全局对象
var G Application

// Plugin 插件接口
type Plugin interface {
	Name() string
	Init() error
	Fini() error
}

// Flagger 实现该接口的插件可通过fs设置命令行参数
type Flagger interface {
	SetFlags(fs *FlagSet)
}

// Runner 每一个实现该接口的插件都会有一个后台goroutine运行该接口.
// 当程序结束时, 会调用ctx的cancel来通知插件程序结束运行.
type Runner interface {
	Run(ctx context.Context) error
}

// Configurator 配置器接口
type Configurator interface {
	ToString(configs map[string]interface{}) string
	WriteToFile(filename string, configs map[string]interface{}) error
	LoadFromFile(filename string, configs map[string]interface{}) error
}

// Options 框架定义的命令行选项
type Options struct {
	Version       bool
	ConfigFile    string
	ConfigExample string
}

// Application 应用框架
type Application struct {
	// 命令行FlagSet, 默认使用flag.CommandLine
	CommandLine *flag.FlagSet

	// 框架命令行选项
	Options Options

	// 配置器, 默认使用jsonc.Configurator
	Configurator Configurator

	// 版本信息函数, 为nil则无法输出版本信息
	VersionInfo func() string

	// 配置
	configs map[string]interface{}

	// 插件列表
	plugins []Plugin
}

// ConfigInfo 配置信息
func (app *Application) ConfigInfo() string {
	if app.Configurator == nil {
		app.Configurator = jsonc.Configurator
	}
	return app.Configurator.ToString(app.configs)
}

// Register 注册插件
func (app *Application) Register(p Plugin, c interface{}) {
	for _, v := range app.plugins {
		if v.Name() == p.Name() {
			panic(fmt.Sprintf("%q plugin is registered", p.Name()))
		}
	}
	app.plugins = append(app.plugins, p)

	if app.configs == nil {
		app.configs = make(map[string]interface{})
	}
	if c != nil {
		app.configs[p.Name()] = c
	}
}

// Main 运行主程序
func (app *Application) Main(args []string) {
	var err error

	// 解析命令行参数
	if err = app.parseCommandLine(args); err != nil {
		fmt.Fprintf(os.Stderr, "parse command line: %v\n", err)
		exit(3)
	}

	// 执行命令
	if err = app.doCommand(); err != nil {
		fmt.Fprintf(os.Stderr, "do command: %v\n", err)
		exit(3)
	}

	// 加载配置
	if err = app.loadConfig(); err != nil {
		fmt.Fprintf(os.Stderr, "load config: %v\n", err)
		exit(3)
	}

	// 输出程序基本信息
	app.printInfo()

	// 运行插件
	if err = app.runPlugins(); err != nil {
		fmt.Fprintf(os.Stderr, "run plugins: %v\n", err)
		exit(3)
	}
}

func (app *Application) setupCommandLine() {
	if app.VersionInfo != nil {
		app.CommandLine.BoolVar(&app.Options.Version, "version", app.Options.Version, "输出版本信息")
	}
	if len(app.configs) > 0 {
		app.CommandLine.StringVar(&app.Options.ConfigFile, "config-file", app.Options.ConfigFile, "指定配置文件")
		app.CommandLine.StringVar(&app.Options.ConfigExample, "config-example", app.Options.ConfigExample, "生成配置示例")
	}

	for _, p := range app.plugins {
		if f, ok := p.(Flagger); ok {
			f.SetFlags(&FlagSet{
				prefix:  p.Name() + ".",
				flagset: app.CommandLine,
			})
		}
	}
}

func (app *Application) parseCommandLine(args []string) error {
	if app.CommandLine == nil {
		app.CommandLine = flag.CommandLine
	}
	app.setupCommandLine()
	return app.CommandLine.Parse(args[1:])
}

func (app *Application) doCommand() (err error) {
	var quit bool

	if app.Options.Version {
		if app.VersionInfo == nil {
			return errors.New("Application.VersionInfo is nil")
		}
		fmt.Fprintf(os.Stdout, "%s", app.VersionInfo())
		quit = true
	} else if app.Options.ConfigExample != "" {
		if app.Configurator == nil {
			app.Configurator = jsonc.Configurator
		}
		if err = app.Configurator.WriteToFile(app.Options.ConfigExample, app.configs); err != nil {
			return fmt.Errorf("generate config example: %v", err)
		}
		fmt.Fprintf(os.Stdout, "generate config example %s success\n", app.Options.ConfigExample)
		quit = true
	}

	if quit {
		exit(0)
	}

	return nil
}

func (app *Application) loadConfig() (err error) {
	if app.Options.ConfigFile == "" {
		return nil
	}
	if len(app.configs) <= 0 {
		return nil
	}
	if app.Configurator == nil {
		app.Configurator = jsonc.Configurator
	}
	return app.Configurator.LoadFromFile(app.Options.ConfigFile, app.configs)
}

func (app *Application) printInfo() {
	if app.VersionInfo != nil {
		log.Infof("version info:\n%s", app.VersionInfo())
	}
	if len(app.configs) > 0 {
		log.Infof("config info:\n%s", app.ConfigInfo())
	}
}

func (app *Application) runPlugins() (err error) {
	// init plugins
	for _, p := range app.plugins {
		if err = p.Init(); err != nil {
			log.Errorw("init", "plugin", p.Name(), "error", err)
			return fmt.Errorf("init %s plugin: %v", p.Name(), err)
		}
		log.Debugw("init success", "plugin", p.Name())
	}

	// quit signal
	quit := make(chan error, 1)
	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		ch := make(chan os.Signal, 1)
		signal.Notify(ch, os.Interrupt, os.Kill)
		<-ch
		cancel()
		time.Sleep(10 * time.Second)
		quit <- errors.New("wait 10s, force exit")
	}()

	// run plugins
	var wg sync.WaitGroup
	for _, p := range app.plugins {
		if r, ok := p.(Runner); ok {
			wg.Add(1)
			go func(n string, r Runner) {
				if err := r.Run(ctx); err != nil {
					quit <- fmt.Errorf("run %s plugin: %v", n, err)
				}
				wg.Done()
			}(p.Name(), r)
		}
	}
	go func() {
		wg.Wait()
		quit <- nil
	}()

	// fini plugins
	rerr := <-quit
	for i := len(app.plugins) - 1; i >= 0; i-- {
		p := app.plugins[i]
		if err = p.Fini(); err != nil {
			log.Errorw("fini", "plugin", p.Name(), "error", err)
			continue
		}
		log.Debugw("fini success", "plugin", p.Name())
	}
	time.Sleep(100 * time.Millisecond)

	return rerr
}

package main

import (
	"fmt"
	"os"

	"git.ablecloud.cn/ablecloud/ac-comm-lib/pluginapp"
	"git.ablecloud.cn/ablecloud/ac-comm-lib/pluginapp/configurator/tomlc"
	_ "git.ablecloud.cn/ablecloud/ac-comm-lib/pluginapp/plugins/backend"
)

type Options struct {
	Development bool
}

type Config struct {
	Environment string
	Node        string
	Addr        string
}

type Plugin struct {
	Options Options
	Config  Config
}

func (p *Plugin) Name() string {
	return "main"
}

func (p *Plugin) SetFlags(fs *pluginapp.FlagSet) {
	fs.BoolVar(&p.Options.Development, "development", false, "启动开发模式")
}

func (p *Plugin) Init() error {
	fmt.Println(p.Options.Development, p.Config.Environment, p.Config.Node, p.Config.Addr)
	return nil
}

func (p *Plugin) Fini() error {
	return nil
}

var G = &Plugin{
	Config: Config{
		Environment: "production",
		Node:        "node1",
		Addr:        "localhost:8000",
	},
}

func init() {
	pluginapp.G.Register(G, &G.Config)
}

func main() {
	pluginapp.G.VersionInfo = func() string { return "v0.0.1\n" }
	pluginapp.G.Configurator = tomlc.Configurator
	pluginapp.G.Main(os.Args)
}

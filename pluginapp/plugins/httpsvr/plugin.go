package httpsvr

import (
	"context"
	"net/http"

	"git.ablecloud.cn/ablecloud/ac-comm-lib/httputils"
	"git.ablecloud.cn/ablecloud/ac-comm-lib/pluginapp"
)

type Config struct {
	Addr       string
	LogVerbose int64
}

type Plugin struct {
	Config  Config
	Handler http.Handler
}

func (p *Plugin) Name() string {
	return "httpsvr"
}

func (p *Plugin) Init() error {
	return nil
}

func (p *Plugin) Fini() error {
	return nil
}

func (p *Plugin) Run(ctx context.Context) error {
	h := httputils.NewVerboseHandler(httputils.NewVerbose(p.Config.LogVerbose), nil, p.Handler)
	svr := http.Server{Addr: p.Config.Addr, Handler: h}
	go func() {
		<-ctx.Done()
		svr.Close()
	}()
	if err := svr.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}

var G = Plugin{
	Config: Config{
		Addr: ":10000",
	},
}

func init() {
	pluginapp.G.Register(&G, &G.Config)
}

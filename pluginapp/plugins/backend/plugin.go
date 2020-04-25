package backend

import (
	"context"
	"net/http"

	"git.ablecloud.cn/ablecloud/ac-comm-lib/pluginapp"
)

type Config struct {
	Addr string
}

type Plugin struct {
	Config Config
	http.ServeMux
}

func (p *Plugin) Name() string {
	return "backend"
}

func (p *Plugin) Init() error {
	return nil
}

func (p *Plugin) Fini() error {
	return nil
}

func (p *Plugin) Run(ctx context.Context) error {
	svr := http.Server{Addr: p.Config.Addr, Handler: &p.ServeMux}
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
		Addr: ":6060",
	},
}

func init() {
	pluginapp.G.Register(&G, &G.Config)
}

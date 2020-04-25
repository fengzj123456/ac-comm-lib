package pluginapp_test

import (
	"testing"

	"github.com/ironzhang/x-pearls/log"

	"git.ablecloud.cn/ablecloud/ac-comm-lib/pluginapp"
	"git.ablecloud.cn/ablecloud/ac-comm-lib/pluginapp/plugins/backend"
)

func TestMain(m *testing.M) {
	log.SetLevel("debug")
	m.Run()
}

func TestPlugins(t *testing.T) {
	var _ = backend.G

	args := []string{"test"}
	pluginapp.G.Main(args)
}

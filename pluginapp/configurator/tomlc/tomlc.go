package tomlc

import (
	"bytes"
	"fmt"

	"github.com/BurntSushi/toml"
	"github.com/ironzhang/x-pearls/config/tomlcfg"
)

type configurator struct{}

func (configurator) ToString(configs map[string]interface{}) string {
	var b bytes.Buffer
	enc := toml.NewEncoder(&b)
	enc.Indent = "\t"
	enc.Encode(configs)
	return b.String()
}

func (configurator) WriteToFile(filename string, configs map[string]interface{}) error {
	return tomlcfg.TOML.WriteToFile(filename, configs)
}

func (configurator) LoadFromFile(filename string, configs map[string]interface{}) (err error) {
	var m map[string]toml.Primitive
	if err = tomlcfg.TOML.LoadFromFile(filename, &m); err != nil {
		return err
	}
	for k, prim := range m {
		if cfg, ok := configs[k]; ok {
			if err = toml.PrimitiveDecode(prim, cfg); err != nil {
				return fmt.Errorf("unmarshal value of %q: %v", k, err)
			}
		}
	}
	return nil
}

var Configurator configurator

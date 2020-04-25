package jsonc

import (
	"encoding/json"
	"fmt"

	"github.com/ironzhang/x-pearls/config"
)

type configurator struct{}

func (configurator) ToString(configs map[string]interface{}) string {
	data, _ := json.MarshalIndent(configs, "", "\t")
	return string(data)
}

func (configurator) WriteToFile(filename string, configs map[string]interface{}) error {
	return config.JSON.WriteToFile(filename, configs)
}

func (configurator) LoadFromFile(filename string, configs map[string]interface{}) (err error) {
	var m map[string]json.RawMessage
	if err = config.JSON.LoadFromFile(filename, &m); err != nil {
		return err
	}
	for k, raw := range m {
		if cfg, ok := configs[k]; ok {
			if err = json.Unmarshal(raw, cfg); err != nil {
				return fmt.Errorf("unmarshal value of %q: %v", k, err)
			}
		}
	}
	return nil
}

var Configurator configurator

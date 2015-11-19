package config_parser

import (
	"encoding/json"

	debian_config "github.com/bborbe/debian/config"
)

type ConfigParser interface {
	ParseConfig(content []byte) (*debian_config.Config, error)
}

type configParser struct {
}

func New() *configParser {
	return new(configParser)
}

func (c *configParser) ParseConfig(content []byte) (*debian_config.Config, error) {
	config := debian_config.DefaultConfig()
	if err := json.Unmarshal(content, config); err != nil {
		return nil, err
	}
	return config, nil
}

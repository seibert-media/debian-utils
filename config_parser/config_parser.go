package config_parser

import (
	"encoding/json"

	"io/ioutil"

	debian_config "github.com/bborbe/debian_utils/config"
	io_util "github.com/bborbe/io/util"
)

type ConfigParser interface {
	ParseFileToConfig(config *debian_config.Config, path string) (*debian_config.Config, error)
	ParseContentToConfig(config *debian_config.Config, content []byte) (*debian_config.Config, error)
}

type configParser struct {
}

func New() *configParser {
	return new(configParser)
}

func (c *configParser) ParseContentToConfig(config *debian_config.Config, content []byte) (*debian_config.Config, error) {
	if err := json.Unmarshal(content, config); err != nil {
		return nil, err
	}
	return config, nil
}

func (c *configParser) ParseFileToConfig(config *debian_config.Config, path string) (*debian_config.Config, error) {
	var content []byte
	var err error
	if path, err = io_util.NormalizePath(path); err != nil {
		return nil, err
	}
	if content, err = ioutil.ReadFile(path); err != nil {
		return nil, err
	}
	return c.ParseContentToConfig(config, content)
}

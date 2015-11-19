package config_builder

import (
	"fmt"

	debian_config "github.com/bborbe/debian/config"
	"github.com/bborbe/log"
)

type ConfigBuilder interface {
	Name(name string) error
	Version(version string) error
	Section(section string) error
	Priority(priority string) error
	Architecture(architecture string) error
	Maintainer(maintainer string) error
	Description(description string) error
	AddFile(source string, target string) error
	Build() *debian_config.Config
}

var logger = log.DefaultLogger

type configBuilder struct {
	config *debian_config.Config
}

func New() *configBuilder {
	c := new(configBuilder)
	c.config = debian_config.DefaultConfig()
	return c
}

func (c *configBuilder) Section(section string) error {
	logger.Debugf("Section %s", section)
	if len(section) == 0 {
		return fmt.Errorf("section empty")
	}
	c.config.Section = section
	return nil
}

func (c *configBuilder) Priority(priority string) error {
	logger.Debugf("Priority %s", priority)
	if len(priority) == 0 {
		return fmt.Errorf("priority empty")
	}
	c.config.Priority = priority
	return nil
}

func (c *configBuilder) Architecture(architecture string) error {
	logger.Debugf("Architecture %s", architecture)
	if len(architecture) == 0 {
		return fmt.Errorf("architecture empty")
	}
	c.config.Architecture = architecture
	return nil
}

func (c *configBuilder) Maintainer(maintainer string) error {
	logger.Debugf("Maintainer %s", maintainer)
	if len(maintainer) == 0 {
		return fmt.Errorf("maintainer empty")
	}
	c.config.Maintainer = maintainer
	return nil
}

func (c *configBuilder) Description(description string) error {
	logger.Debugf("Description %s", description)
	if len(description) == 0 {
		return fmt.Errorf("description empty")
	}
	c.config.Description = description
	return nil
}

func (c *configBuilder) Name(name string) error {
	logger.Debugf("Name %s", name)
	if len(name) == 0 {
		return fmt.Errorf("name empty")
	}
	c.config.Name = name
	return nil
}

func (c *configBuilder) Version(version string) error {
	logger.Debugf("Version %s", version)
	if len(version) == 0 {
		return fmt.Errorf("version empty")
	}
	c.config.Version = version
	return nil
}

func (c *configBuilder) AddFile(source string, target string) error {
	if len(source) == 0 {
		return fmt.Errorf("source empty")
	}
	if len(target) == 0 {
		return fmt.Errorf("target empty")
	}
	c.config.Files = append(c.config.Files, debian_config.File{Source: source, Target: target})
	return nil
}

func (c *configBuilder) Build() *debian_config.Config {
	return c.config
}

package config_builder

import (
	"fmt"

	debian_config "github.com/seibert-media/debian-utils/config"
	"github.com/golang/glog"
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
	AddDepend(value string) error
	AddConflict(value string) error
	AddProvide(value string) error
	AddReplace(value string) error
	Build() *debian_config.Config
}

type configBuilder struct {
	config *debian_config.Config
}

func New() *configBuilder {
	return NewWithConfig(debian_config.DefaultConfig())
}

func NewWithConfig(config *debian_config.Config) *configBuilder {
	c := new(configBuilder)
	c.config = config
	return c
}

func (c *configBuilder) Section(section string) error {
	glog.V(2).Infof("Section %s", section)
	if len(section) == 0 {
		return fmt.Errorf("section empty")
	}
	c.config.Section = section
	return nil
}

func (c *configBuilder) Priority(priority string) error {
	glog.V(2).Infof("Priority %s", priority)
	if len(priority) == 0 {
		return fmt.Errorf("priority empty")
	}
	c.config.Priority = priority
	return nil
}

func (c *configBuilder) Architecture(architecture string) error {
	glog.V(2).Infof("Architecture %s", architecture)
	if len(architecture) == 0 {
		return fmt.Errorf("architecture empty")
	}
	c.config.Architecture = architecture
	return nil
}

func (c *configBuilder) Maintainer(maintainer string) error {
	glog.V(2).Infof("Maintainer %s", maintainer)
	if len(maintainer) == 0 {
		return fmt.Errorf("maintainer empty")
	}
	c.config.Maintainer = maintainer
	return nil
}

func (c *configBuilder) Description(description string) error {
	glog.V(2).Infof("Description %s", description)
	if len(description) == 0 {
		return fmt.Errorf("description empty")
	}
	c.config.Description = description
	return nil
}

func (c *configBuilder) Name(name string) error {
	glog.V(2).Infof("Name %s", name)
	if len(name) == 0 {
		return fmt.Errorf("name empty")
	}
	c.config.Name = name
	return nil
}

func (c *configBuilder) Version(version string) error {
	glog.V(2).Infof("Version %s", version)
	if len(version) == 0 {
		return fmt.Errorf("version empty")
	}
	c.config.Version = version
	return nil
}

func (c *configBuilder) AddDepend(value string) error {
	if len(value) == 0 {
		return fmt.Errorf("depend empty")
	}
	c.config.Depends = append(c.config.Depends, value)
	return nil
}

func (c *configBuilder) AddConflict(value string) error {
	if len(value) == 0 {
		return fmt.Errorf("depend empty")
	}
	c.config.Conflicts = append(c.config.Conflicts, value)
	return nil
}

func (c *configBuilder) AddProvide(value string) error {
	if len(value) == 0 {
		return fmt.Errorf("depend empty")
	}
	c.config.Provides = append(c.config.Provides, value)
	return nil
}

func (c *configBuilder) AddReplace(value string) error {
	if len(value) == 0 {
		return fmt.Errorf("depend empty")
	}
	c.config.Replaces = append(c.config.Replaces, value)
	return nil
}

func (c *configBuilder) AddFile(source string, target string) error {
	if len(source) == 0 {
		return fmt.Errorf("add file failed. source is empty")
	}
	if len(target) == 0 {
		return fmt.Errorf("add file failed. target is empty")
	}
	c.config.Files = append(c.config.Files, debian_config.File{Source: source, Target: target})
	return nil
}

func (c *configBuilder) Build() *debian_config.Config {
	return c.config
}

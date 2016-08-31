package config

type File struct {
	Source  string `json:"source"`
	Target  string `json:"target"`
	Extract bool   `json:"extract"`
}

type Config struct {
	Name         string   `json:"name"`
	Version      string   `json:"version"`
	Files        []File   `json:"files"`
	Section      string   `json:"section"`
	Priority     string   `json:"priority"`
	Architecture string   `json:"architecture"`
	Maintainer   string   `json:"maintainer"`
	Description  string   `json:"description"`
	Homepage     string   `json:"homepage"`
	Depends      []string `json:"depends"`
	Conflicts    []string `json:"conflicts"`
	Provides     []string `json:"provides"`
	Replaces     []string `json:"replaces"`
	Postrm       string   `json:"postrm"`
	Postinst     string   `json:"postinst"`
	Prerm        string   `json:"prerm"`
	Preinst      string   `json:"preinst"`
}

func DefaultConfig() *Config {
	c := new(Config)
	c.Section = "base"
	c.Priority = "optional"
	c.Architecture = "all"
	c.Maintainer = "Benjamin Borbe <bborbe@rocketnews.de>"
	c.Description = "-"
	c.Files = []File{}
	c.Depends = []string{}
	c.Conflicts = []string{}
	c.Provides = []string{}
	c.Replaces = []string{}
	return c
}

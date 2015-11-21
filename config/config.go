package config

type File struct {
	Source string `json:"source"`
	Target string `json:"target"`
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
	return c
}

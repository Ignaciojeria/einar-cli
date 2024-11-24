package domain

type EinarCli struct {
	Project       string         `json:"project"`
	Template      Template       `json:"template"`
	Installations []Installation `json:"installations"`
	Components    []Component    `json:"components"`
}

type Template struct {
	Tag string `json:"tag"`
	URL string `json:"url"`
}

type Installation struct {
	Name      string   `json:"name"`
	Unique    string   `json:"unique"`
	Libraries []string `json:"libraries"`
}

type Component struct {
	Kind string `json:"kind"`
	Name string `json:"name"`
}

func (c EinarCli) IsInstalled(component string) bool {
	for _, installation := range c.Installations {
		if installation.Name == component {
			return true
		}
	}
	return false
}

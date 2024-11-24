package domain

type EinarTemplate struct {
	BaseTemplate         BaseTemplate          `json:"base_template"`
	InstallationsBase    []InstallationsBase   `json:"installations_base"`
	InstallationCommands []InstallationCommand `json:"installation_commands"`
	ComponentCommands    []ComponentCommands   `json:"component_commands"`
}

type BaseTemplate struct {
	Description string       `json:"description"`
	Folders     []BaseFolder `json:"folders"`
	Files       []BaseFile   `json:"files"`
}

type BaseFolder struct {
	SourceDir      string `json:"source_dir"`
	DestinationDir string `json:"destination_dir"`
}

type BaseFile struct {
	SourceFile      string `json:"source_file"`
	DestinationFile string `json:"destination_file"`
}

type InstallationFolder struct {
	SourceDir      string `json:"source_dir"`
	DestinationDir string `json:"destination_dir"`
	IocDiscovery   bool   `json:"ioc_discovery"`
}

type InstallationFile struct {
	SourceFile     string          `json:"source_file"`
	DestinationDir string          `json:"destination_dir"`
	IocDiscovery   bool            `json:"ioc_discovery"`
	Port           Port            `json:"port"`
	ReplaceHolders []ReplaceHolder `json:"replace_holders"`
}

type InstallationCommand struct {
	Name           string               `json:"name"`
	Unique         string               `json:"unique"`
	SourceDir      string               `json:"source_dir"`
	DestinationDir string               `json:"destination_dir"`
	Folders        []InstallationFolder `json:"folders"`
	Files          []InstallationFile   `json:"files"`
	Command        string               `json:"command"`
	Libraries      []string             `json:"libraries"`
	DependsOn      []string             `json:"depends_on"`
}

type InstallationsBase struct {
	Name    string `json:"name"`
	Library string `json:"library"`
}

type ComponentCommands struct {
	Kind           string          `json:"kind"`
	Name           string          `json:"name"`
	ComponentFiles []ComponentFile `json:"files"`
	DependsOn      []string        `json:"depends_on"`
}

type ComponentFile struct {
	SourceFile          string               `json:"source_file"`
	DestinationDir      string               `json:"destination_dir"`
	IocDiscovery        bool                 `json:"ioc_discovery"`
	HasComponentDir     bool                 `json:"has_component_dir"`
	AppendAtStart       string               `json:"append_at_start"`
	AppendAtEnd         string               `json:"append_at_end"`
	Port                Port                 `json:"port"`
	ReplaceHolders      []ReplaceHolder      `json:"replace_holders"`
	LiteralReplacements []LiteralReplacement `json:"literal_replacements"`
}

type Port struct {
	SourceFile     string `json:"source_file"`
	DestinationDir string `json:"destination_dir"`
}

type ReplaceHolder struct {
	Kind          string `json:"kind"`
	Name          string `json:"name"`
	AppendAtStart string `json:"append_at_start"`
	AppendAtEnd   string `json:"append_at_end"`
}

type LiteralReplacement struct {
	Target      string `json:"target"`
	Replacement string `json:"replacement"`
}

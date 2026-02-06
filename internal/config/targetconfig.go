package config

type TargetConfig struct {
	TemplateName string         `yaml:"template_name"`
	Variables    map[string]any `yaml:"variables"`
	Scripts      map[string]any `yaml:"scripts"`
	Definitions  map[string]any `yaml:"definitions"`
}

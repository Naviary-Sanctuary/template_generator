package config

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/pelletier/go-toml"
)

const (
	DefaultConfigFile  = "tg.config.toml"
	DefaultTemplateDir = ".tg"
	TemplateConfigFile = "template.toml"
)

type Config struct {
	TemplatesDir string         `toml:"templates_dir"`
	Defaults     map[string]any `toml: "defaults,omitempty"`
}

type Metadata struct {
	Name        string `toml:"name"`
	Description string `toml:"description,omitempty"`
	Author      string `toml:"author,omitempty"`
}

type Template struct {
	Metadata  Metadata            `toml:"metadata"`
	Variables map[string]Variable `toml:"variables"`
	Rules     Rules               `toml:"rules"`
	Version   string              `toml:"version,omitempty"`
}

type Variable struct {
	Default     any    `toml:"default,omitempty"`
	Description string `toml:"description,omitempty"`
	Type        string `toml:"type, omitempty"`
}

type Rules struct {
	Ignores  []string          `toml:"ignores,omitempty"`
	Includes []string          `toml:"includes,omitempty"`
	Renames  map[string]string `toml:"renames,omitempty"`
}

func NewConfig() *Config {
	return &Config{
		TemplatesDir: DefaultTemplateDir,
		Defaults:     make(map[string]any),
	}
}

func Load(path string) (*Config, error) {
	if path == "" {
		path = DefaultConfigFile
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := toml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	if config.TemplatesDir == "" {
		config.TemplatesDir = DefaultTemplateDir
	}

	return &config, nil
}

func (config *Config) Save(path string) error {
	if path == "" {
		path = DefaultConfigFile
	}

	data, err := toml.Marshal(config)
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	if err := os.WriteFile(path, data, 0644); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func LoadTemplate(dir string) (*Template, error) {
	configPath := filepath.Join(dir, TemplateConfigFile)

	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read template config file: %w", err)
	}

	var template Template
	if err := toml.Unmarshal(data, &template); err != nil {
		return nil, fmt.Errorf("failed to parse template config file: %w", err)
	}

	if template.Metadata.Name == "" {
		template.Metadata.Name = filepath.Base(dir)
	}

	if template.Variables == nil {
		template.Variables = make(map[string]Variable)
	}

	for name, variable := range template.Variables {
		if variable.Type == "" {
			variable.Type = "any"
			template.Variables[name] = variable
		}
	}

	return &template, nil
}

func SaveTemplate(dir string, tmpl *Template) error {
	configPath := filepath.Join(dir, TemplateConfigFile)

	// Create directory if it doesn't exist
	if err := os.MkdirAll(dir, 0755); err != nil {
		return fmt.Errorf("failed to create template directory: %w", err)
	}

	data, err := toml.Marshal(tmpl)
	if err != nil {
		return fmt.Errorf("failed to marshal template: %w", err)
	}

	if err := os.WriteFile(configPath, data, 0644); err != nil {
		return fmt.Errorf("failed to write template config: %w", err)
	}

	return nil
}

func (config *Config) Validate() error {
	if config.TemplatesDir == "" {
		return fmt.Errorf("templates_dir cannot be empty")
	}

	if info, err := os.Stat(config.TemplatesDir); err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("templates directory '%s' does not exist", config.TemplatesDir)
		}
		return fmt.Errorf("failed to check templates directory: %w", err)
	} else if !info.IsDir() {
		return fmt.Errorf("'%s' is not a directory", config.TemplatesDir)
	}

	return nil
}

func (t *Template) Validate() error {
	if t.Metadata.Name == "" {
		return fmt.Errorf("template name is required")
	}

	for name, variable := range t.Variables {
		if err := validateVariable(name, variable); err != nil {
			return err
		}
	}

	// Check for conflicting rules
	if len(t.Rules.Includes) > 0 && len(t.Rules.Ignores) > 0 {
		// This is allowed, but includes take precedence
		// Just log a warning in verbose mode
	}

	return nil
}

func validateVariable(name string, v Variable) error {
	validTypes := []string{"string", "number", "boolean", "array"}
	if v.Type != "" {
		valid := false
		for _, t := range validTypes {
			if v.Type == t {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("variable '%s': invalid type '%s'", name, v.Type)
		}
	}

	if v.Default != nil && v.Type != "" {
		if err := validateValueType(v.Default, v.Type); err != nil {
			return fmt.Errorf("variable '%s': default value error: %w", name, err)
		}
	}

	return nil
}

func validateValueType(value interface{}, expectedType string) error {
	switch expectedType {
	case "string":
		if _, ok := value.(string); !ok {
			return fmt.Errorf("expected string, got %T", value)
		}
	case "number":
		switch value.(type) {
		case int, int64, float64:
		default:
			return fmt.Errorf("expected number, got %T", value)
		}
	case "boolean":
		if _, ok := value.(bool); !ok {
			return fmt.Errorf("expected boolean, got %T", value)
		}
	case "array":
		switch value.(type) {
		case []interface{}, []string:
		default:
			return fmt.Errorf("expected array, got %T", value)
		}
	}
	return nil
}

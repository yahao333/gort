package config

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Version      string                 `yaml:"version"`
	Environments map[string]Environment `yaml:"environments"`
	Providers    map[string]Provider    `yaml:"providers"`
	Defaults     map[string]interface{} `yaml:"defaults"`
}

type Environment struct {
	Provider  string                 `yaml:"provider"`
	Region    string                 `yaml:"region"`
	Variables map[string]interface{} `yaml:"variables"`
	Tags      map[string]string      `yaml:"tags"`
	Backend   *Backend               `yaml:"backend,omitempty"`
}

type Provider struct {
	Type       string                 `yaml:"type"`
	Version    string                 `yaml:"version"`
	Properties map[string]interface{} `yaml:"properties"`
}

type Backend struct {
	Type   string                 `yaml:"type"`
	Config map[string]interface{} `yaml:"config"`
}

// LoadConfig loads configuration from the specified path
func LoadConfig(path string) (*Config, error) {
	if path == "" {
		path = "gort.yaml"
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config: %w", err)
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid configuration: %w", err)
	}

	return &config, nil
}

func (c *Config) Validate() error {
	if len(c.Environments) == 0 {
		return fmt.Errorf("no environments defined")
	}

	for name, env := range c.Environments {
		if env.Provider == "" {
			return fmt.Errorf("provider not specified for environment %s", name)
		}

		if _, exists := c.Providers[env.Provider]; !exists {
			return fmt.Errorf("undefined provider '%s' referenced in environment %s",
				env.Provider, name)
		}
	}

	return nil
}

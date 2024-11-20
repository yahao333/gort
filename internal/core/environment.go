package core

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/yahao333/gort/internal/logging"
	"gopkg.in/yaml.v3"
)

type EnvironmentManager struct {
	configPath string
	lockPath   string
	mu         sync.Mutex
	logger     *logging.Logger
}

type EnvironmentConfig struct {
	Variables map[string]string         `yaml:"variables"`
	Secrets   map[string]string         `yaml:"secrets"`
	Providers map[string]ProviderConfig `yaml:"providers"`
}

type ProviderConfig struct {
	Type       string                 `yaml:"type"`
	Properties map[string]interface{} `yaml:"properties"`
}

func NewEnvironmentManager(configPath string, logger *logging.Logger) *EnvironmentManager {
	return &EnvironmentManager{
		configPath: configPath,
		lockPath:   filepath.Join(filepath.Dir(configPath), "locks"),
		logger:     logger,
	}
}

func (em *EnvironmentManager) LoadEnvironment(name string) (*EnvironmentConfig, error) {
	configFile := filepath.Join(em.configPath, fmt.Sprintf("%s.yaml", name))
	data, err := os.ReadFile(configFile)
	if err != nil {
		return nil, fmt.Errorf("failed to read environment config: %w", err)
	}

	var config EnvironmentConfig
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse environment config: %w", err)
	}

	return &config, nil
}

func (em *EnvironmentManager) LockEnvironment(name string) error {
	em.mu.Lock()
	defer em.mu.Unlock()

	lockFile := filepath.Join(em.lockPath, fmt.Sprintf("%s.lock", name))
	if err := os.MkdirAll(em.lockPath, 0755); err != nil {
		return fmt.Errorf("failed to create lock directory: %w", err)
	}

	if _, err := os.Stat(lockFile); err == nil {
		return fmt.Errorf("environment %s is already locked", name)
	}

	return os.WriteFile(lockFile, []byte{}, 0644)
}

func (em *EnvironmentManager) UnlockEnvironment(name string) error {
	em.mu.Lock()
	defer em.mu.Unlock()

	lockFile := filepath.Join(em.lockPath, fmt.Sprintf("%s.lock", name))
	return os.Remove(lockFile)
}

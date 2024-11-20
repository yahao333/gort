package plugin

import (
	"fmt"
	"path/filepath"
	"plugin"
)

type PluginManager struct {
	pluginDir string
	plugins   map[string]interface{}
	logger    *Logger
}

type PluginMetadata struct {
	Name        string   `json:"name"`
	Version     string   `json:"version"`
	Author      string   `json:"author"`
	Description string   `json:"description"`
	Interfaces  []string `json:"interfaces"`
}

func NewPluginManager(pluginDir string, logger *Logger) *PluginManager {
	return &PluginManager{
		pluginDir: pluginDir,
		plugins:   make(map[string]interface{}),
		logger:    logger,
	}
}

func (pm *PluginManager) LoadPlugin(name string) error {
	pluginPath := filepath.Join(pm.pluginDir, fmt.Sprintf("%s.so", name))

	p, err := plugin.Open(pluginPath)
	if err != nil {
		return fmt.Errorf("failed to open plugin: %w", err)
	}

	// Load metadata
	metadataSymbol, err := p.Lookup("Metadata")
	if err != nil {
		return fmt.Errorf("plugin metadata not found: %w", err)
	}

	metadata, ok := metadataSymbol.(*PluginMetadata)
	if !ok {
		return fmt.Errorf("invalid plugin metadata type")
	}

	// Load plugin instance
	newSymbol, err := p.Lookup("New")
	if err != nil {
		return fmt.Errorf("plugin constructor not found: %w", err)
	}

	constructor, ok := newSymbol.(func() interface{})
	if !ok {
		return fmt.Errorf("invalid plugin constructor type")
	}

	instance := constructor()
	pm.plugins[metadata.Name] = instance
	pm.logger.Infof("Loaded plugin: %s v%s", metadata.Name, metadata.Version)

	return nil
}

func (pm *PluginManager) GetPlugin(name string) (interface{}, error) {
	plugin, exists := pm.plugins[name]
	if !exists {
		return nil, fmt.Errorf("plugin %s not found", name)
	}
	return plugin, nil
}

package plugin

import (
	"context"
	"fmt"
	"path/filepath"
	"plugin"
	"sync"
)

// PluginManager handles plugin lifecycle and management
type PluginManager struct {
	mu          sync.RWMutex
	pluginDir   string
	plugins     map[string]*PluginInfo
	logger      *Logger
	initialized bool
}

// PluginInfo stores plugin metadata and instance
type PluginInfo struct {
	Metadata *PluginMetadata
	Instance Plugin
	Path     string
	Loaded   bool
}

// PluginMetadata contains plugin metadata
type PluginMetadata struct {
	Name        string            `json:"name"`
	Version     string            `json:"version"`
	Type        string            `json:"type"`
	Author      string            `json:"author"`
	Description string            `json:"description"`
	Properties  map[string]string `json:"properties"`
}

// NewPluginManager creates a new plugin manager
func NewPluginManager(pluginDir string, logger *Logger) *PluginManager {
	return &PluginManager{
		pluginDir: pluginDir,
		plugins:   make(map[string]*PluginInfo),
		logger:    logger,
	}
}

// Initialize scans plugin directory and loads plugin metadata
func (pm *PluginManager) Initialize(ctx context.Context) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	if pm.initialized {
		return nil
	}

	// Scan plugin directory
	files, err := filepath.Glob(filepath.Join(pm.pluginDir, "*.so"))
	if err != nil {
		return fmt.Errorf("failed to scan plugin directory: %w", err)
	}

	for _, file := range files {
		if err := pm.loadPluginMetadata(file); err != nil {
			pm.logger.Errorf("Failed to load plugin metadata from %s: %v", file, err)
			continue
		}
	}

	pm.initialized = true
	return nil
}

// LoadPlugin loads a specific plugin
func (pm *PluginManager) LoadPlugin(ctx context.Context, name string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	info, exists := pm.plugins[name]
	if !exists {
		return fmt.Errorf("plugin %s not found", name)
	}

	if info.Loaded {
		return nil
	}

	p, err := plugin.Open(info.Path)
	if err != nil {
		return fmt.Errorf("failed to open plugin: %w", err)
	}

	// Load plugin instance
	newSymbol, err := p.Lookup("New")
	if err != nil {
		return fmt.Errorf("plugin constructor not found: %w", err)
	}

	constructor, ok := newSymbol.(func() Plugin)
	if !ok {
		return fmt.Errorf("invalid plugin constructor type")
	}

	instance := constructor()
	info.Instance = instance
	info.Loaded = true

	// Initialize plugin
	if err := instance.Init(nil); err != nil {
		return fmt.Errorf("failed to initialize plugin: %w", err)
	}

	pm.logger.Infof("Loaded plugin: %s v%s", info.Metadata.Name, info.Metadata.Version)
	return nil
}

// UnloadPlugin unloads a specific plugin
func (pm *PluginManager) UnloadPlugin(ctx context.Context, name string) error {
	pm.mu.Lock()
	defer pm.mu.Unlock()

	info, exists := pm.plugins[name]
	if !exists {
		return fmt.Errorf("plugin %s not found", name)
	}

	if !info.Loaded {
		return nil
	}

	if err := info.Instance.Shutdown(ctx); err != nil {
		return fmt.Errorf("failed to shutdown plugin: %w", err)
	}

	info.Instance = nil
	info.Loaded = false
	return nil
}

// GetPlugin returns a loaded plugin instance
func (pm *PluginManager) GetPlugin(name string) (Plugin, error) {
	pm.mu.RLock()
	defer pm.mu.RUnlock()

	info, exists := pm.plugins[name]
	if !exists {
		return nil, fmt.Errorf("plugin %s not found", name)
	}

	if !info.Loaded {
		return nil, fmt.Errorf("plugin %s is not loaded", name)
	}

	return info.Instance, nil
}

// loadPluginMetadata loads plugin metadata without fully loading the plugin
func (pm *PluginManager) loadPluginMetadata(path string) error {
	p, err := plugin.Open(path)
	if err != nil {
		return fmt.Errorf("failed to open plugin: %w", err)
	}

	metadataSymbol, err := p.Lookup("Metadata")
	if err != nil {
		return fmt.Errorf("plugin metadata not found: %w", err)
	}

	metadata, ok := metadataSymbol.(*PluginMetadata)
	if !ok {
		return fmt.Errorf("invalid plugin metadata type")
	}

	pm.plugins[metadata.Name] = &PluginInfo{
		Metadata: metadata,
		Path:     path,
		Loaded:   false,
	}

	return nil
}

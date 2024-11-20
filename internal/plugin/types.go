package plugin

import "context"

// Plugin represents the base interface that all plugins must implement
type Plugin interface {
	// Init initializes the plugin with configuration
	Init(config map[string]interface{}) error
	// Name returns the plugin name
	Name() string
	// Version returns the plugin version
	Version() string
	// Shutdown cleanly shuts down the plugin
	Shutdown(ctx context.Context) error
}

// ProviderPlugin represents a infrastructure provider plugin
type ProviderPlugin interface {
	Plugin
	// CreateResource creates a new resource
	CreateResource(ctx context.Context, spec ResourceSpec) (*Resource, error)
	// DeleteResource deletes an existing resource
	DeleteResource(ctx context.Context, id string) error
	// UpdateResource updates an existing resource
	UpdateResource(ctx context.Context, id string, spec ResourceSpec) (*Resource, error)
	// GetResource gets resource details
	GetResource(ctx context.Context, id string) (*Resource, error)
}

// ResourceSpec defines the specification for a resource
type ResourceSpec struct {
	Type       string                 `json:"type"`
	Name       string                 `json:"name"`
	Properties map[string]interface{} `json:"properties"`
}

// Resource represents a managed resource
type Resource struct {
	ID         string                 `json:"id"`
	Type       string                 `json:"type"`
	Name       string                 `json:"name"`
	Properties map[string]interface{} `json:"properties"`
	Status     string                 `json:"status"`
}

// HookPlugin represents a plugin that can hook into various lifecycle events
type HookPlugin interface {
	Plugin
	// PreCreate is called before resource creation
	PreCreate(ctx context.Context, spec ResourceSpec) error
	// PostCreate is called after resource creation
	PostCreate(ctx context.Context, resource Resource) error
	// PreDelete is called before resource deletion
	PreDelete(ctx context.Context, id string) error
	// PostDelete is called after resource deletion
	PostDelete(ctx context.Context, id string) error
}

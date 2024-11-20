package core

import (
	"fmt"
	"time"

	"github.com/yahao333/gort/internal/logging"
	"github.com/yahao333/gort/internal/state"
)

type ResourceType string

const (
	ResourceTypeInstance ResourceType = "instance"
	ResourceTypeDatabase ResourceType = "database"
	ResourceTypeNetwork  ResourceType = "network"
)

type ResourceState string

const (
	ResourceStatePending  ResourceState = "pending"
	ResourceStateCreating ResourceState = "creating"
	ResourceStateRunning  ResourceState = "running"
	ResourceStateFailed   ResourceState = "failed"
	ResourceStateDeleted  ResourceState = "deleted"
)

type ResourceManager struct {
	stateManager *state.StateManager
	logger       *logging.Logger
}

type ResourceSpec struct {
	Name         string                 `json:"name"`
	Type         ResourceType           `json:"type"`
	Provider     string                 `json:"provider"`
	Properties   map[string]interface{} `json:"properties"`
	Dependencies []string               `json:"dependencies"`
}

type ResourceStatus struct {
	State       ResourceState `json:"state"`
	Message     string        `json:"message"`
	LastUpdated time.Time     `json:"last_updated"`
}

func NewResourceManager(stateManager *state.StateManager, logger *logging.Logger) *ResourceManager {
	return &ResourceManager{
		stateManager: stateManager,
		logger:       logger,
	}
}

func (rm *ResourceManager) CreateResource(env string, spec *ResourceSpec) error {
	rm.logger.Infof("Creating resource %s of type %s", spec.Name, spec.Type)

	// Check dependencies
	for _, dep := range spec.Dependencies {
		if err := rm.checkDependency(env, dep); err != nil {
			return fmt.Errorf("dependency check failed: %w", err)
		}
	}

	// Create resource
	status := &ResourceStatus{
		State:       ResourceStateCreating,
		LastUpdated: time.Now(),
	}

	if err := rm.updateResourceStatus(env, spec.Name, status); err != nil {
		return fmt.Errorf("failed to update resource status: %w", err)
	}

	return nil
}

func (rm *ResourceManager) DeleteResource(env string, name string) error {
	rm.logger.Infof("Deleting resource %s from environment %s", name, env)

	status := &ResourceStatus{
		State:       ResourceStateDeleted,
		LastUpdated: time.Now(),
	}

	return rm.updateResourceStatus(env, name, status)
}

func (rm *ResourceManager) checkDependency(env string, resourceName string) error {
	status, err := rm.getResourceStatus(env, resourceName)
	if err != nil {
		return err
	}

	if status.State != ResourceStateRunning {
		return fmt.Errorf("dependency %s is not in running state", resourceName)
	}

	return nil
}

func (rm *ResourceManager) updateResourceStatus(env string, name string, status *ResourceStatus) error {
	state, err := rm.stateManager.LoadState(env)
	if err != nil {
		return err
	}

	if state.Resources == nil {
		state.Resources = make(map[string]interface{})
	}

	state.Resources[name] = status
	return rm.stateManager.SaveState(env, state)
}

func (rm *ResourceManager) getResourceStatus(env string, name string) (*ResourceStatus, error) {
	state, err := rm.stateManager.LoadState(env)
	if err != nil {
		return nil, err
	}

	if status, ok := state.Resources[name].(*ResourceStatus); ok {
		return status, nil
	}

	return nil, fmt.Errorf("resource %s not found", name)
}

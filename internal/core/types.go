package core

// Environment represents a deployment environment
type Environment struct {
	Name      string
	Provider  string
	Region    string
	Variables map[string]interface{}
}

// Deployment represents a deployment state
type Deployment struct {
	ID          string
	Environment string
	Version     string
	Status      string
	Resources   []Resource
}

// Resource represents an infrastructure resource
type Resource struct {
	Type       string
	Name       string
	Properties map[string]interface{}
}

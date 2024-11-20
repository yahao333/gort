package provider

// Provider defines the interface for infrastructure providers
type Provider interface {
	// Initialize sets up the provider
	Initialize() error

	// Plan returns the planned changes
	Plan(env string) (*PlanResult, error)

	// Apply applies the planned changes
	Apply(plan *PlanResult) error
}

type PlanResult struct {
	Changes     []Change
	AddCount    int
	UpdateCount int
	DeleteCount int
}

type Change struct {
	Type     string // add, update, delete
	Resource string
	Before   interface{}
	After    interface{}
}

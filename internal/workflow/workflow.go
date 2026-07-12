package workflow

// Debug controls whether detailed debug logs are printed.
var Debug bool

// Service holds shared dependencies for workflows.
type Service struct {
	// Add workflow dependencies here as needed
}

// New builds a new workflow service container.
func New() *Service {
	return &Service{}
}

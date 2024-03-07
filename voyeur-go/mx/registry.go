package mx

type Registry struct {
	root ManagedResource
}

var (
	registry Registry
)

func InitRegistry(root ManagedResource) {
	registry.root = root
}

func Health() HealthState {
	if registry.root != nil {
		return registry.root.Health()
	}
	return Unknown
}

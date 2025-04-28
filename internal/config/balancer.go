package config

// Balancer represents config for the balancer.
type Balancer struct {
	// Backends is a list of <host>:<port> strings.
	Backends []string `yaml:"backends"`
	// Port to listen.
	Port uint `yaml:"port"`
}

// DefaultForBalancer returns default config for balancer.
func DefaultForBalancer() Balancer {
	return Balancer{
		Backends: []string{},
		Port:     8080,
	}
}

package config

// Balancer represents config for the balancer.
type Balancer struct {
	// Backends is a list of <host>:<port> strings.
	Backends []string `yaml:"backends"`
	// Port to listen.
	Port uint32 `yaml:"port"`
	// Strategy name to use. Available are:
	//	- RoundRobin.
	Strategy string `yaml:"strategy"`
	// Healthcheck config.
	Heathcheck Heathcheck `yaml:"healthcheck"`
}

// Heathcheck represents config for healhchecks.
type Heathcheck struct {
	// CheckTimeoutSeconds is period between checking backend's health.
	CheckTimeoutSeconds uint32 `yaml:"check_timeout_seconds"`
	// RequestTimeoutSeconds is timeout for check health request.
	RequestTimeoutSeconds uint32 `yaml:"request_timeout_seconds"`
}

// DefaultForBalancer returns default config for balancer.
func DefaultForBalancer() Balancer {
	return Balancer{
		Backends: []string{},
		Port:     8080,
		Strategy: "RoundRobin",
		Heathcheck: Heathcheck{
			CheckTimeoutSeconds:   60,
			RequestTimeoutSeconds: 30,
		},
	}
}

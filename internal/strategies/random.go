package strategies

import (
	"math/rand"
	"sync"
)

// Random strategy for balancing requests to backends.
type Random struct {
	backendAvailable map[string]*backendState
	backends         []string
}

// NewRandom creates new Random strategy.
func NewRandom(backends []string) *Random {
	backendAvailable := make(map[string]*backendState, len(backends))
	for _, backend := range backends {
		backendAvailable[backend] = &backendState{
			rwLock:    &sync.RWMutex{},
			available: false,
		}
	}

	return &Random{
		backendAvailable: backendAvailable,
		backends:         backends,
	}
}

// ChooseBackend returns backend host which is ready to receive request.
func (r *Random) ChooseBackend() string {
	order := rand.Perm(len(r.backends))

	for _, i := range order {
		candidate := r.backends[i]
		state := r.backendAvailable[candidate]

		state.rwLock.RLock()
		available := state.available
		state.rwLock.RUnlock()

		if available {
			return candidate
		}
	}

	return ""
}

// UpdateBackendHealth marks given backend health.
func (r *Random) UpdateBackendHealth(backend string, healthy bool) {
	updateBackendHeath(r.backendAvailable, backend, healthy)
}

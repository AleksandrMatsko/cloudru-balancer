package strategies

import (
	"sync"
)

type backendState struct {
	rwLock    *sync.RWMutex
	available bool
}

// RoundRobin is a cyclic balancer strategy.
type RoundRobin struct {
	backendAvailable map[string]*backendState
	allBackends      []string
	indexLock        sync.Locker
	startIndex       int
}

// NewRoundRobin creates RoundRobin.
func NewRoundRobin(backends []string) *RoundRobin {
	availabilityMap := make(map[string]*backendState)
	for i := range backends {
		availabilityMap[backends[i]] = &backendState{
			rwLock:    &sync.RWMutex{},
			available: false,
		}
	}

	return &RoundRobin{
		backendAvailable: availabilityMap,
		allBackends:      backends,
		indexLock:        &sync.Mutex{},
		startIndex:       0,
	}
}

// ChooseBackend returns backend host which is ready to receive request.
func (rr *RoundRobin) ChooseBackend() string {
	rr.indexLock.Lock()
	startIndex := rr.startIndex
	rr.startIndex += 1
	rr.indexLock.Unlock()

	for i := range rr.allBackends {
		candidate := rr.allBackends[(startIndex+i)%len(rr.allBackends)]
		state := rr.backendAvailable[candidate]

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
func (rr *RoundRobin) UpdateBackendHealth(backend string, healthy bool) {
	updateBackendHeath(rr.backendAvailable, backend, healthy)
}

func updateBackendHeath(backendAvailable map[string]*backendState, backend string, healthy bool) {
	if state, ok := backendAvailable[backend]; ok {
		state.rwLock.RLock()
		same := state.available == healthy
		state.rwLock.RUnlock()

		if same {
			return
		}

		state.rwLock.Lock()
		state.available = healthy
		state.rwLock.Unlock()
	}
}

package strategies

import (
	"sync"
)

// RoundRobin is a cyclic balancer strategy.
type RoundRobin struct {
	backendAvailable map[string]bool
	// allBackends is a slice of available backend's hosts. Methods must not mutate the slice.
	allBackends []string
	// locker protects startIndex and backendAvailable map
	locker     sync.Locker
	startIndex int
}

// NewRoundRobin creates RoundRobin.
func NewRoundRobin(backends []string) *RoundRobin {
	availabilityMap := make(map[string]bool)
	for i := range backends {
		availabilityMap[backends[i]] = false
	}

	return &RoundRobin{
		backendAvailable: availabilityMap,
		allBackends:      backends,
		locker:           &sync.Mutex{},
		startIndex:       0,
	}
}

// ChooseBackend returns backend host which is ready to receive request.
func (rr *RoundRobin) ChooseBackend() string {
	rr.locker.Lock()
	defer rr.locker.Unlock()

	for i := range rr.allBackends {
		candidate := rr.allBackends[(rr.startIndex+i)%len(rr.allBackends)]
		if rr.backendAvailable[candidate] {
			rr.startIndex = (rr.startIndex + i + 1) % len(rr.allBackends)
			return candidate
		}
	}

	return ""
}

// UpdateBackendHealth
func (rr *RoundRobin) UpdateBackendHealth(backend string, healthy bool) {
	rr.locker.Lock()
	defer rr.locker.Unlock()

	if _, ok := rr.backendAvailable[backend]; ok {
		rr.backendAvailable[backend] = healthy
	}
}

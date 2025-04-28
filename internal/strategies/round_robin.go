package strategies

import (
	"sync"
)

// RoundRobin is a cyclic balancer strategy.
type RoundRobin struct {
	// allBackends is a slice of available backend's hosts. Methods must not mutate the slice.
	allBackends []string
	// locker protects startIndex
	locker     sync.Locker
	startIndex int
}

// NewRoundRobin creates RoundRobin.
func NewRoundRobin(backends []string) *RoundRobin {
	if len(backends) <= 1 {
		return &RoundRobin{
			allBackends: backends,
		}
	}

	return &RoundRobin{
		allBackends: backends,
		locker:      &sync.Mutex{},
		startIndex:  0,
	}
}

// SendOrder returns ordered backends according to RoundRobin strategy.
func (rr *RoundRobin) SendOrder() []string {
	orderedBackends := make([]string, len(rr.allBackends))
	copy(orderedBackends, rr.allBackends)

	if len(rr.allBackends) <= 1 {
		return orderedBackends
	}

	rr.locker.Lock()
	startIndexCopy := rr.startIndex
	rr.startIndex = (rr.startIndex + 1) % len(rr.allBackends)
	rr.locker.Unlock()

	for i := range orderedBackends {
		orderedBackends[i] = rr.allBackends[(startIndexCopy+i)%len(orderedBackends)]
	}

	return orderedBackends
}

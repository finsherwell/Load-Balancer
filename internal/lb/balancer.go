package lb

import (
	"sync"
	"time"
)

// Algorithm represents a load balancing algorithm
type LoadBalancerAlgorithm interface {
	SelectBackend(backends []*Backend) *Backend
}

// Balancer manages load balancing
type Balancer struct {
	backends       []*Backend
	mux            sync.Mutex
	algorithm      LoadBalancerAlgorithm
	lastSwitchTime time.Time
}

// NewBalancer creates a new load balancer
func NewBalancer(backends []*Backend, algorithm LoadBalancerAlgorithm) *Balancer {
	return &Balancer{
		backends:       backends,
		algorithm:      algorithm,
		lastSwitchTime: time.Now(),
	}
}

// SwitchAlgorithm changes the load balancing algorithm
func (b *Balancer) SwitchAlgorithm(algorithm LoadBalancerAlgorithm) {
	b.mux.Lock()
	defer b.mux.Unlock()
	b.algorithm = algorithm
}

// SelectBackend selects a backend based on the current algorithm
func (b *Balancer) SelectBackend() *Backend {
	b.mux.Lock()
	defer b.mux.Unlock()
	return b.algorithm.SelectBackend(b.backends)
}

// Monitor and switch algorithm
func (b *Balancer) MonitorAndSwitch() {
	b.mux.Lock()
	defer b.mux.Unlock()

	var totalConnections int32

	for _, backends := range b.backends {
		totalConnections += backends.GetConnections()
	}

	// Switch algorithm based on load conditions
	if totalConnections > 50 { // Example threshold
		if _, ok := b.algorithm.(*LeastConnections); !ok {
			b.SwitchAlgorithm(&LeastConnections{})
		}
	} else if totalConnections < 20 { // Example threshold
		if _, ok := b.algorithm.(*RoundRobin); !ok {
			b.SwitchAlgorithm(&RoundRobin{})
		}
	} else {
		if _, ok := b.algorithm.(*Random); !ok {
			b.SwitchAlgorithm(&Random{})
		}
	}

	// Optionally, add a cooldown period to avoid constant switching
	if time.Since(b.lastSwitchTime) > time.Minute {
		b.lastSwitchTime = time.Now()
	}
}

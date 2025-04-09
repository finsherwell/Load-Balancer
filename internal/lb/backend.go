package lb

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"sync"
	"sync/atomic"
	"time"
)

// Backend represents a backend server
type Backend struct {
	URL            *url.URL
	Alive          bool
	mux            sync.RWMutex
	ReverseProxy   *httputil.ReverseProxy
	MaxConnections int
	Connections    int32
	Healthy        bool
}

// NewBackend creates a new backend server
func NewBackend(urlStr string, maxConn int) (*Backend, error) {
	url, err := url.Parse(urlStr)
	if err != nil {
		return nil, err
	}

	b := &Backend{
		URL:            url,
		Alive:          true,
		ReverseProxy:   httputil.NewSingleHostReverseProxy(url),
		MaxConnections: maxConn,
		Connections:    0,
		Healthy:        true,
	}

	return b, nil
}

// SetAlive sets the status of the backend to alive or not
func (b *Backend) SetAlive(alive bool) {
	b.mux.Lock()
	defer b.mux.Unlock()
	b.Alive = alive
}

// IsAlive checks if the backend is alive or not
func (b *Backend) IsAlive() bool {
	b.mux.RLock()
	defer b.mux.RUnlock()
	return b.Alive
}

// IncrementConnections increments the active connections count
func (b *Backend) IncrementConnections() int32 {
	return atomic.AddInt32(&b.Connections, 1)
}

// DecrementConnections decrements the active connections count
func (b *Backend) DecrementConnections() int32 {
	return atomic.AddInt32(&b.Connections, -1)
}

// GetConnections returns the current connections count
func (b *Backend) GetConnections() int32 {
	return atomic.LoadInt32(&b.Connections)
}

// SetHealthy sets the healthy status of the backend
func (b *Backend) SetHealthy(health bool) {
	b.mux.Lock()
	defer b.mux.Unlock()
	b.Healthy = health
}

// IsHealthy returns the healthy status of the backend
func (b *Backend) IsHealthy() bool {
	b.mux.RLock()
	defer b.mux.RUnlock()
	return b.Healthy
}

// CanAcceptConnection checks if backend can accept more connections
func (b *Backend) CanAcceptConnection() bool {
	return b.GetConnections() < int32(b.MaxConnections) && b.IsAlive() && b.IsHealthy()
}

// HealthCheck pings the backend to check its health
func (b *Backend) HealthCheck() bool {
	client := http.Client{Timeout: 5 * time.Second}
	resp, err := client.Head(b.URL.String() + "/health")
	isHealthy := err == nil && resp.StatusCode == http.StatusOK
	b.SetHealthy(isHealthy)
	return isHealthy
}

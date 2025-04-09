package server

import (
	"fmt"
	"log"
	"net/http"

	"github.com/finsherwell/Load-Balancer/internal/config"
	"github.com/finsherwell/Load-Balancer/internal/lb"
)

// StartServer initializes the HTTP server and load balancer
func StartServer(configPath string) error {
	// Load config
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	// Initialize backends
	var backends []*lb.Backend
	for _, urlStr := range cfg.Backends {
		backend, err := lb.NewBackend(urlStr, cfg.MaxConns)
		if err != nil {
			return fmt.Errorf("failed to initialize backend: %w", err)
		}
		backends = append(backends, backend)
	}

	// Initialize load balancer with the selected algorithm
	var algorithm lb.LoadBalancerAlgorithm
	switch cfg.Algorithm {
	case "round_robin":
		algorithm = &lb.RoundRobin{}
	case "least_connections":
		algorithm = &lb.LeastConnections{}
	case "random":
		algorithm = &lb.Random{}
	default:
		return fmt.Errorf("unsupported algorithm: %s", cfg.Algorithm)
	}

	// Create the balancer
	balancer := lb.NewBalancer(backends, algorithm)

	// HTTP handler to forward requests
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		backend := balancer.SelectBackend()
		if backend == nil {
			http.Error(w, "No available backends", http.StatusServiceUnavailable)
			return
		}

		// Forward the request to the selected backend
		backend.ReverseProxy.ServeHTTP(w, r)
	})

	// Start server
	log.Printf("Starting server on :8080")
	return http.ListenAndServe(":8080", nil)
}

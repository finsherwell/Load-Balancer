package lb

import (
	"math/rand"
	"time"
)

type Random struct{}

func (r *Random) SelectBackend(backends []*Backend) *Backend {
	rand.NewSource(time.Now().UnixNano())
	return backends[rand.Intn(len(backends))]
}

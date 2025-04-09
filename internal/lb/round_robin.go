package lb

type RoundRobin struct {
	current int
}

func (rr *RoundRobin) SelectBackend(backends []*Backend) *Backend {
	if len(backends) == 0 {
		return nil
	}
	backend := backends[rr.current]
	rr.current = (rr.current + 1) % len(backends)
	return backend
}

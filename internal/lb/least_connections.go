package lb

type LeastConnections struct{}

func (lc *LeastConnections) SelectBackend(backends []*Backend) *Backend {
	var selected *Backend
	for _, backend := range backends {
		if selected == nil || backend.GetConnections() < selected.GetConnections() {
			selected = backend
		}
	}
	return selected
}

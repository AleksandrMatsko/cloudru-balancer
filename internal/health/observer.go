package health

// Observer should be used for tracking backends' healph.
type Observer interface {
	// UpdateBackendHealth for given backend.
	UpdateBackendHealth(backend string, heathy bool)
}

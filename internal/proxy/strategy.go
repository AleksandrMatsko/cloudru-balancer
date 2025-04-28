package proxy

// Strategy is the interface used to decide which backend server use.
type Strategy interface {
	// ChooseBackend returns backend host which is ready to receive request.
	ChooseBackend() string
}

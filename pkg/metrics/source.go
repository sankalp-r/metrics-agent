package metrics

// Source interface represents different metrics source
type Source interface {
	// Collect metrics
	Collect() ([]byte, error)
}

package multiplexer

// Multiplexer is a generic interface for different multiplexing mechanisms.
type Multiplexer interface {
	// Add adds a file descriptor (socket, file, etc.) to be monitored for readability.
	Add(fd int) error

	// Remove removes a file descriptor from monitoring.
	Remove(fd int) error

	// Wait waits for events on monitored file descriptors and returns a list of readable file descriptors.
	Wait() ([]int, error)

	// Close releases any resources associated with the multiplexer.
	Close() error
}

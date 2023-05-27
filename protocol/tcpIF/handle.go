package tcpIF

import (
	"context"
	"net"
)

type Handler interface {
	// Handle client connection
	Handle(ct context.Context, conn net.Conn)

	// Close handler
	Close() error
}

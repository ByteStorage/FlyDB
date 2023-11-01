package multiplexer

import (
	"errors"
	"golang.org/x/net/ipv6"
	"net"
)

// avportMultiplexer is an implementation of the Multiplexer interface using avport.
type avportMultiplexer struct {
	conn       *ipv6.PacketConn
	fds        map[int]bool
	addChan    chan int
	removeChan chan int
	stopChan   chan struct{}
}

func NewAvportMultiplexer() (Multiplexer, error) {
	conn, err := net.ListenPacket("udp6", "[::1]:0")
	if err != nil {
		return nil, err
	}

	pc := ipv6.NewPacketConn(conn)

	return &avportMultiplexer{
		conn:       pc,
		fds:        make(map[int]bool),
		addChan:    make(chan int),
		removeChan: make(chan int),
		stopChan:   make(chan struct{}),
	}, nil
}

func (a *avportMultiplexer) Add(fd int) error {
	a.addChan <- fd
	return nil
}

func (a *avportMultiplexer) Remove(fd int) error {
	a.removeChan <- fd
	return nil
}

func (a *avportMultiplexer) Wait() ([]int, error) {
	readable := make([]int, 0)
	controlBuf := make([]byte, 1024)

	for {
		select {
		case fd := <-a.addChan:
			if !a.fds[fd] {
				a.fds[fd] = true
				// Configure socket options for avport (e.g., multicast group, flow control)
				// Set up appropriate socket options for your application
			}
		case fd := <-a.removeChan:
			if a.fds[fd] {
				delete(a.fds, fd)
				// Remove any socket options for avport
			}
		case <-a.stopChan:
			return nil, errors.New("multiplexer closed")
		default:
			n, _, _, err := a.conn.ReadFrom(controlBuf)
			if err != nil {
				return nil, err
			}
			// Process control data to determine readable file descriptors
			for i := 0; i < n; i++ {
				readable = append(readable, int(controlBuf[i]))
			}
			if len(readable) > 0 {
				return readable, nil
			}
		}
	}
}

func (a *avportMultiplexer) Close() error {
	err := a.conn.Close()
	if err != nil {
		return err
	}
	close(a.stopChan)
	return nil
}

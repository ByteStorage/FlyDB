package multiplexer

import (
	"errors"
	"golang.org/x/sys/unix"
)

// epollMultiplexer is an implementation of the Multiplexer interface using epoll.
type epollMultiplexer struct {
	epfd       int
	events     []unix.EpollEvent
	fds        map[int]bool
	addChan    chan int
	removeChan chan int
	stopChan   chan struct{}
}

func NewEpollMultiplexer() (Multiplexer, error) {
	epfd, err := unix.EpollCreate1(0)
	if err != nil {
		return nil, err
	}

	events := make([]unix.EpollEvent, 64) // Adjust the size as needed

	return &epollMultiplexer{
		epfd:       epfd,
		events:     events,
		fds:        make(map[int]bool),
		addChan:    make(chan int),
		removeChan: make(chan int),
		stopChan:   make(chan struct{}),
	}, nil
}

func (e *epollMultiplexer) Add(fd int) error {
	e.addChan <- fd
	return nil
}

func (e *epollMultiplexer) Remove(fd int) error {
	e.removeChan <- fd
	return nil
}

func (e *epollMultiplexer) Wait() ([]int, error) {
	readable := make([]int, 0)

	select {
	case fd := <-e.addChan:
		if !e.fds[fd] {
			e.fds[fd] = true
			event := unix.EpollEvent{
				Events: unix.EPOLLIN,
				Fd:     int32(fd),
			}
			unix.EpollCtl(e.epfd, unix.EPOLL_CTL_ADD, fd, &event)
		}
	case fd := <-e.removeChan:
		if e.fds[fd] {
			delete(e.fds, fd)
			unix.EpollCtl(e.epfd, unix.EPOLL_CTL_DEL, fd, nil)
		}
	case <-e.stopChan:
		return nil, errors.New("multiplexer closed")
	default:
		n, err := unix.EpollWait(e.epfd, e.events, -1)
		if err != nil {
			return nil, err
		}

		for i := 0; i < n; i++ {
			readable = append(readable, int(e.events[i].Fd))
		}
	}

	return readable, nil
}

func (e *epollMultiplexer) Close() error {
	unix.Close(e.epfd)
	close(e.stopChan)
	return nil
}

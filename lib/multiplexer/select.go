package multiplexer

import (
	"errors"
)

// SelectMultiplexer is an implementation of the Multiplexer interface using select.
type selectMultiplexer struct {
	fds        map[int]bool
	readable   chan int
	addChan    chan int
	removeChan chan int
	stopChan   chan struct{}
}

func NewSelectMultiplexer() (Multiplexer, error) {
	return &selectMultiplexer{
		fds:        make(map[int]bool),
		readable:   make(chan int),
		addChan:    make(chan int),
		removeChan: make(chan int),
		stopChan:   make(chan struct{}),
	}, nil
}

func (s *selectMultiplexer) Add(fd int) error {
	s.addChan <- fd
	return nil
}

func (s *selectMultiplexer) Remove(fd int) error {
	s.removeChan <- fd
	return nil
}

func (s *selectMultiplexer) Wait() ([]int, error) {
	readable := make([]int, 0)

	select {
	case fd := <-s.readable:
		readable = append(readable, fd)
	case fd := <-s.addChan:
		s.fds[fd] = true
	case fd := <-s.removeChan:
		delete(s.fds, fd)
	case <-s.stopChan:
		return nil, errors.New("multiplexer closed")
	}

	return readable, nil
}

func (s *selectMultiplexer) Close() error {
	close(s.stopChan)
	return nil
}

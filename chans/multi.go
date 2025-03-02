package chans

import (
	"sync"
)

type MultiChan[T any] struct {
	src       chan T
	listeners map[chan T]bool

	mu sync.Mutex
}

func (multi *MultiChan[T]) readMsg(msg T) {
	multi.mu.Lock()
	defer multi.mu.Unlock()

	for listener, ok := range multi.listeners {
		if !ok {
			continue
		}

		listener <- msg
	}
}

func (multi *MultiChan[T]) listen() {
	for msg := range multi.src {
		multi.readMsg(msg)
	}
}

func (multi *MultiChan[T]) Chan() chan<- T {
	return multi.src
}

func (multi *MultiChan[T]) Send(data T) {
	multi.src <- data
}

func (multi *MultiChan[T]) Register() <-chan T {
	multi.mu.Lock()
	defer multi.mu.Unlock()

	listener := make(chan T)
	multi.listeners[listener] = true

	return listener
}

func (multi *MultiChan[T]) Unregister(src <-chan T) {
	multi.mu.Lock()
	defer multi.mu.Unlock()

	for listener, ok := range multi.listeners {
		if listener == src && ok {
			delete(multi.listeners, listener)
			close(listener)
		}
	}
}

func (multi *MultiChan[T]) Close() {
	close(multi.src)

	multi.mu.Lock()
	defer multi.mu.Unlock()

	for listener, ok := range multi.listeners {
		if ok {
			delete(multi.listeners, listener)
			close(listener)
		}
	}
}

func NewMultiChan[T any]() *MultiChan[T] {
	multi := &MultiChan[T]{src: make(chan T), listeners: make(map[chan T]bool)}
	go multi.listen()

	return multi
}

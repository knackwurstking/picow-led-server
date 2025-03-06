package event

import (
	"fmt"
	"sync"
)

type (
	Name            string
	Listener[T any] func(a T)
)

type Event[T any] struct {
	Listeners map[Name][]Listener[T]

	mutex *sync.Mutex
}

func NewEvent[T any]() *Event[T] {
	return &Event[T]{
		Listeners: make(map[Name][]Listener[T]),
		mutex:     &sync.Mutex{},
	}
}

func (e *Event[T]) Dispatch(n Name, d T) {
	defer e.mutex.Unlock()
	e.mutex.Lock()

	if e.Listeners[n] == nil {
		return
	}

	wg := sync.WaitGroup{}

	for _, listener := range e.Listeners[n] {
		wg.Add(1)
		go func() {
			defer wg.Done()
			go listener(d)
		}()
	}

	wg.Wait()
}

func (e *Event[T]) On(n Name, l Listener[T]) {
	defer e.mutex.Unlock()
	e.mutex.Lock()

	if e.Listeners[n] == nil {
		e.Listeners[n] = []Listener[T]{l}
		return
	}

	lSig := fmt.Sprintf("%v", l)
	for _, l2 := range e.Listeners[n] {
		if fmt.Sprintf("%v", l2) == lSig {
			return
		}
	}

	e.Listeners[n] = append(e.Listeners[n], l)
}

func (e *Event[T]) Off(n Name, l Listener[T]) {
	defer e.mutex.Unlock()
	e.mutex.Lock()

	if e.Listeners[n] == nil {
		return
	}

	listeners := make([]Listener[T], 0)

	lSig := fmt.Sprintf("%v", l)
	for _, listener := range e.Listeners[n] {
		if fmt.Sprintf("%v", listener) == lSig {
			listeners = append(listeners, listener)
		}
	}

	e.Listeners[n] = listeners
}

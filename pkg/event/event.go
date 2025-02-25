package event

type (
	Name            string
	Listener[T any] func(a T)
)

type Event[T any] struct {
	// TODO: Store listeners in a map
	// TODO: Add Off method to remove a listener
}

func NewEvent[T any]() *Event[T] {
	return &Event[T]{}
}

func (e *Event[T]) Dispatch() {}

func (e *Event[T]) On(n Name, l Listener[T]) *Listener[T] {
	// ...

	return &l // for off method
}

func (e *Event[T]) Off(n Name, l *Listener[T]) {
	// ...
}

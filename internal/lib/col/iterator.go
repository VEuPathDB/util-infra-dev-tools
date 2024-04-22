package col

type Iterator[T any] interface {
	HasNext() bool
	Next() T
}

type MapIterator[K comparable, V any] interface {
	HasNext() bool
	Next() (K, V)
}

package col

import "github.com/sirupsen/logrus"

type OrderedMap[K comparable, V any] interface {
	Contains(key K) bool
	Put(key K, val V)
	Get(key K) (V, bool)
	Require(key K) V
	Delete(key K) bool
	Size() int
	ForEach(fn func(K, V))
	Iterator() MapIterator[K, V]
}

func NewOrderedMap[K comparable, V any](size int) OrderedMap[K, V] {
	return &orderedMap[K, V]{
		order: make([]K, 0, size),
		index: make(map[K]V, size),
	}
}

type orderedMap[K comparable, V any] struct {
	order []K
	index map[K]V
}

func (o orderedMap[K, V]) Contains(key K) bool {
	_, ok := o.index[key]
	return ok
}

func (o *orderedMap[K, V]) Put(key K, value V) {
	a := !o.Contains(key)
	o.index[key] = value
	if a {
		o.order = append(o.order, key)
	}
}

func (o orderedMap[K, V]) Get(key K) (V, bool) {
	value, exists := o.index[key]
	return value, exists
}

func (o orderedMap[K, V]) Require(key K) V {
	value, exists := o.index[key]

	if !exists {
		logrus.Fatalf("map missing required key: %s", key)
	}

	return value
}

func (o *orderedMap[K, V]) Delete(key K) bool {
	if o.Contains(key) {
		delete(o.index, key)
		o.order, _ = SliceDeleteFirst(key, o.order)
		return true
	}

	return false
}

func (o orderedMap[K, V]) Size() int {
	return len(o.order)
}

func (o orderedMap[K, V]) ForEach(fn func(K, V)) {
	for i := range o.order {
		fn(o.order[i], o.index[o.order[i]])
	}
}

func (o orderedMap[K, V]) Iterator() MapIterator[K, V] {
	return &mapIter[K, V]{o.order, o.index, 0}
}

type mapIter[K comparable, V any] struct {
	o []K
	m map[K]V
	i int
}

func (m mapIter[K, V]) HasNext() bool {
	return m.i < len(m.o)
}

func (m *mapIter[K, V]) Next() (K, V) {
	k := m.i
	m.i++

	return m.o[k], m.m[m.o[k]]
}

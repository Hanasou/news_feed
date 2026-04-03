package ds

import "fmt"

type HeapNode[T fmt.Stringer] struct {
	Value T
}

type Heap[T fmt.Stringer] struct {
	Array []*HeapNode[T]
	Max   bool
}

func MaxHeap[T fmt.Stringer]() *Heap[T] {
	return &Heap[T]{
		[]*HeapNode[T]{},
		true,
	}
}

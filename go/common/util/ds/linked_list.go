package ds

import "fmt"

type ListNode[T fmt.Stringer] struct {
	Value T
	Next  *ListNode[T]
	Prev  *ListNode[T]
}

type LinkedList[T fmt.Stringer] struct {
	Head *ListNode[T]
	Tail *ListNode[T]
}

func (list *LinkedList[T]) InsertFirst(value T) *LinkedList[T] {
	node := &ListNode[T]{Value: value}
	if list.Head == nil {
		list.Head = node
	}
	if list.Tail == nil {
		list.Tail = node
		return list
	}
	head := list.Head
	list.Head = node
	node.Next = head
	head.Prev = node
	return list
}

func (list *LinkedList[T]) InsertLast(value T) *LinkedList[T] {
	node := &ListNode[T]{Value: value}
	if list.Head == nil {
		list.Head = node
	}
	if list.Tail == nil {
		list.Tail = node
		return list
	}
	tail := list.Tail
	tail.Next = node
	node.Prev = tail

	return list
}

func (list *LinkedList[T]) PeekFirst() T {
	return list.Head.Value
}

func (list *LinkedList[T]) PeekLast() T {
	return list.Tail.Value
}

func (list *LinkedList[T]) PopFirst() T {
	if list.Head == nil {
		return *new(T)
	}
	value := list.Head.Value
	list.Head = list.Head.Next
	return value
}

func (list *LinkedList[T]) PopLast() T {
	if list.Tail == nil {
		return *new(T)
	}
	value := list.Tail.Value
	list.Tail = list.Tail.Prev
	return value
}

func (list *LinkedList[T]) String() string {
	listString := ""
	node := list.Head
	for node != nil {
		listString += node.Value.String()
	}

	return listString
}

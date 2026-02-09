package ds

import (
	"cmp"
	"fmt"
)

type treeGeneric interface {
	fmt.Stringer
	cmp.Ordered
}

type TreeNode[T treeGeneric] struct {
	Value  T
	Parent *TreeNode[T]
	Left   *TreeNode[T]
	Right  *TreeNode[T]
}

type BinarySearchTree[T treeGeneric] struct {
	Root *TreeNode[T]
}

func (tree *BinarySearchTree[T]) Search(value T) *TreeNode[T] {
	return search(value, tree.Root)
}

func search[T treeGeneric](value T, root *TreeNode[T]) *TreeNode[T] {
	if root.Value == value {
		return root
	} else if value > root.Value {
		return search(value, root.Right)
	} else {
		return search(value, root.Left)
	}
}

func (tree *BinarySearchTree[T]) Insert(value T) *BinarySearchTree[T] {
	parent := &TreeNode[T]{}
	root := tree.Root
	for root != nil {
		parent = root
		if value > root.Value {
			root = root.Right
		} else {
			root = root.Left
		}
	}

	if parent == nil {
		tree.Root = &TreeNode[T]{
			Value: value,
		}
	} else if value > parent.Value {
		parent.Right = &TreeNode[T]{
			Value:  value,
			Parent: parent,
		}
	} else {
		parent.Left = &TreeNode[T]{
			Value:  value,
			Parent: parent,
		}
	}
	return tree
}

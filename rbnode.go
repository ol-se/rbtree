package rbtree

import (
	"fmt"
	"strings"
)

// RBNode is a node of a red-black tree.
type RBNode[T any] struct {
	Val     T
	left    *RBNode[T]
	right   *RBNode[T]
	parent  *RBNode[T]
	isBlack bool
}

// Next returns the node with the next closest value and true if this node exists.
func (rbn *RBNode[T]) Next() (*RBNode[T], bool) {
	if rbn.right != nil {
		return rbn.right.leftmost(), true
	}

	for rbn.parent != nil && rbn.parent.right == rbn {
		rbn = rbn.parent
	}

	return rbn.parent, rbn.parent != nil
}

// Prev returns the node with the previous closest value and true if this node exists.
func (rbn *RBNode[T]) Prev() (*RBNode[T], bool) {
	if rbn.left != nil {
		return rbn.left.rightmost(), true
	}

	for rbn.parent != nil && rbn.parent.left == rbn {
		rbn = rbn.parent
	}

	return rbn.parent, rbn.parent != nil
}

// clone recursively copies nodes of the red-black tree to a new red-black tree.
func (rbn *RBNode[T]) clone() *RBNode[T] {
	newNode := &RBNode[T]{
		Val:     rbn.Val,
		isBlack: rbn.isBlack,
	}

	if rbn.left != nil {
		newNode.left = rbn.left.clone()
		newNode.left.parent = newNode
	}

	if rbn.right != nil {
		newNode.right = rbn.right.clone()
		newNode.right.parent = newNode
	}

	return newNode
}

// leftSubtreeIsValid checks the validity of the left subtree.
// leftSubtreeIsValid returns the black height of the tree and true if the tree is valid.
func (rbn *RBNode[T]) leftSubtreeIsValid(initialBlackHeight *int, currentBlackHeight int, cmp func(T, T) int) (int, bool) {
	if rbn.left == nil {
		return currentBlackHeight, true
	}

	if rbn.left.parent != rbn || cmp(rbn.Val, rbn.left.Val) <= 0 {
		return 0, false
	}

	return rbn.left.isValid(initialBlackHeight, currentBlackHeight, cmp)
}

// rightSubtreeIsValid checks the validity of the right subtree.
// rightSubtreeIsValid returns the black height of the tree and true if the tree is valid.
func (rbn *RBNode[T]) rightSubtreeIsValid(initialBlackHeight *int, currentBlackHeight int, cmp func(T, T) int) (int, bool) {
	if rbn.right == nil {
		return currentBlackHeight, true
	}

	if rbn.right.parent != rbn || cmp(rbn.Val, rbn.right.Val) >= 0 {
		return 0, false
	}

	return rbn.right.isValid(initialBlackHeight, currentBlackHeight, cmp)
}

// isValid returns the black height of the red-black tree and true if the tree is valid.
func (rbn *RBNode[T]) isValid(initialBlackHeight *int, currentBlackHeight int, cmp func(T, T) int) (int, bool) {
	if rbn.isBlack {
		currentBlackHeight++
	} else if !rbn.parent.isBlack {
		return 0, false
	}

	if rbn.left == nil && rbn.right == nil {
		if *initialBlackHeight == 0 {
			*initialBlackHeight = currentBlackHeight

			return currentBlackHeight, true
		} else if *initialBlackHeight != currentBlackHeight {
			return 0, false
		}
	}

	leftBlackHeight, ok := rbn.leftSubtreeIsValid(initialBlackHeight, currentBlackHeight, cmp)
	if !ok {
		return 0, false
	}

	rightBlackHeight, ok := rbn.rightSubtreeIsValid(initialBlackHeight, currentBlackHeight, cmp)
	if !ok || leftBlackHeight != rightBlackHeight {
		return 0, false
	}

	return max(leftBlackHeight, currentBlackHeight), true
}

// equalTo recursively checks if both trees have the same structure and nodes.
func (rbn *RBNode[T]) equalTo(anotherRBN *RBNode[T], cmp func(T, T) int) bool {
	if anotherRBN == nil {
		return false
	}

	if cmp(rbn.Val, anotherRBN.Val) != 0 || rbn.isBlack != anotherRBN.isBlack {
		return false
	}

	if rbn.left != nil && !rbn.left.equalTo(anotherRBN.left, cmp) {
		return false
	}

	if rbn.left == nil && anotherRBN.left != nil {
		return false
	}

	if rbn.right != nil && !rbn.right.equalTo(anotherRBN.right, cmp) {
		return false
	}

	if rbn.right == nil && anotherRBN.right != nil {
		return false
	}

	return true
}

// insert adds a new value to the red-black tree.
// If the insertion was successful, the newly inserted node and true are returned.
// Otherwise the existent node and false are returned.
func (rbn *RBNode[T]) insert(val T, cmp func(T, T) int) (*RBNode[T], bool) {
	result := cmp(val, rbn.Val)

	switch {
	case result < 0:
		if rbn.left == nil {
			rbn.left = &RBNode[T]{
				Val:    val,
				parent: rbn,
			}

			return rbn.left, true
		}

		return rbn.left.insert(val, cmp)
	case result > 0:
		if rbn.right == nil {
			rbn.right = &RBNode[T]{
				Val:    val,
				parent: rbn,
			}

			return rbn.right, true
		}

		return rbn.right.insert(val, cmp)
	default:
		return rbn, false
	}
}

// recString makes a multi-string depiction of the tree.
// The tree is aligned left-to-right with the root on the left side of the depiction.
func (rbn *RBNode[T]) recString(result *string, counter int) {
	if rbn.right != nil {
		rbn.right.recString(result, counter+1)
	}

	*result += fmt.Sprintln(strings.Repeat(" ", counter), rbn.Val)

	if rbn.left != nil {
		rbn.left.recString(result, counter+1)
	}
}

// find returns the node pointer and true if a node with particular value was found in the red-black tree.
func (rbn *RBNode[T]) find(val T, cmp func(T, T) int) (*RBNode[T], bool) {
	result := cmp(val, rbn.Val)

	switch {
	case result < 0:
		if rbn.left == nil {
			return nil, false
		}

		return rbn.left.find(val, cmp)
	case result > 0:
		if rbn.right == nil {
			return nil, false
		}

		return rbn.right.find(val, cmp)
	default:
		return rbn, true
	}
}

// leftmost returns the pointer to the node with the smallest value.
func (rbn *RBNode[T]) leftmost() *RBNode[T] {
	if rbn.left != nil {
		return rbn.left.leftmost()
	}

	return rbn
}

// rightmost returns the pointer to the node with the biggest value.
func (rbn *RBNode[T]) rightmost() *RBNode[T] {
	if rbn.right != nil {
		return rbn.right.rightmost()
	}

	return rbn
}

// isBlack returns true if the node is black or nil.
func isBlack[T any](rbn *RBNode[T]) bool {
	return rbn == nil || rbn.isBlack
}

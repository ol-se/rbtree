// Package rbtree provides methods to work with generic red-black tree.
package rbtree

import (
	"cmp"
)

// RBTree is a red-black tree. It contains the size and pointers to the first and the last nodes.
// RBTree consists of red and black nodes.
type RBTree[T any] struct {
	root *RBNode[T]
	cmp  func(T, T) int
	// Min is a pointer to the node with the smallest value of the tree.
	Min *RBNode[T]
	// Left is a pointer to the node with the biggest value of the tree.
	Max *RBNode[T]
	// Count is an amount of nodes in the tree.
	Count int
}

// New returns an empty red-black tree.
// cmp is a pointer to the function to compare user-defined types.
//
// cmp returns the result of comparison:
//
//   - result < 0, if first value is smaller;
//   - result > 0, if first value is bigger;
//   - result == 0, if both values are equal.
//
// For ordered primitive types, use NewOrdered.
func New[T any](cmp func(T, T) int) *RBTree[T] {
	return &RBTree[T]{
		cmp: cmp,
	}
}

// NewOrdered returns an empty red-black tree for primitive types ([cmp.Ordered]).
func NewOrdered[T cmp.Ordered]() *RBTree[T] {
	return New(cmp.Compare[T])
}

// Clone copies the red-black tree to a new red-black tree with the same values and structure.
// Clone returns a new red-black tree.
func (rbt *RBTree[T]) Clone() *RBTree[T] {
	if rbt.root == nil {
		return New[T](rbt.cmp)
	}

	tree := &RBTree[T]{
		root:  rbt.root.clone(),
		cmp:   rbt.cmp,
		Count: rbt.Count,
	}

	tree.Min = tree.root.leftmost()
	tree.Max = tree.root.rightmost()

	return tree
}

// IsValid checks if the tree is a valid red-black tree.
func (rbt *RBTree[T]) IsValid() bool {
	if rbt.cmp == nil {
		return false
	}

	if rbt.root == nil {
		return rbt.Min == nil && rbt.Max == nil && rbt.Count == 0
	}

	if rbt.root.parent != nil || !rbt.root.isBlack {
		return false
	}

	blackHeight, count := 0, 0
	_, isValid := rbt.root.isValid(&blackHeight, 0, rbt.cmp)

	if !isValid || rbt.Min != rbt.root.leftmost() || rbt.Max != rbt.root.rightmost() {
		return false
	}

	for i, ok := rbt.Min, true; ok; i, ok = i.Next() {
		count++
	}

	return count == rbt.Count
}

// EqualTo checks if both trees have the same structure and nodes.
func (rbt *RBTree[T]) EqualTo(anotherRBT *RBTree[T]) bool {
	if anotherRBT == nil {
		return false
	}

	if rbt.root == nil && anotherRBT.root == nil {
		return true
	}

	if rbt.root == nil || anotherRBT.root == nil {
		return false
	}

	if rbt.Count != anotherRBT.Count {
		return false
	}

	return rbt.root.equalTo(anotherRBT.root, rbt.cmp)
}

// Insert adds a new value to the red-black tree and fixes the tree afterwards if necessary.
// If the insertion was successful, the newly inserted node and true are returned.
// Otherwise the existent node and false are returned.
func (rbt *RBTree[T]) Insert(val T) (*RBNode[T], bool) {
	if rbt.root == nil {
		rbt.root = &RBNode[T]{
			Val:     val,
			isBlack: true,
		}

		rbt.Min = rbt.root
		rbt.Max = rbt.root

		rbt.Count++

		return rbt.root, true
	}

	insertedNode, ok := rbt.root.insert(val, rbt.cmp)
	if !ok {
		return insertedNode, false
	}

	if rbt.cmp(val, rbt.Min.Val) < 0 {
		rbt.Min = insertedNode
	} else if rbt.cmp(val, rbt.Max.Val) > 0 {
		rbt.Max = insertedNode
	}

	if !insertedNode.parent.isBlack {
		rbt.solveDoubleRed(insertedNode.parent)
	}

	rbt.Count++

	return insertedNode, true
}

func (rbt *RBTree[T]) String() string {
	if rbt.root == nil {
		return ""
	}

	var result string

	rbt.root.recString(&result, 0)

	return result
}

// Find returns the node pointer and true if a node with particular value was found in the red-black tree.
func (rbt *RBTree[T]) Find(val T) (*RBNode[T], bool) {
	if rbt.root == nil {
		return nil, false
	}

	return rbt.root.find(val, rbt.cmp)
}

// Delete deletes a node with particular value from the red-black tree and fixes the tree if necessary.
// Delete returns the deleted value and true if deletion was successful. It returns an empty value and false otherwise.
func (rbt *RBTree[T]) Delete(val T) (T, bool) {
	var del T

	if rbt.root == nil {
		return del, false
	}

	rbnDelete, ok := rbt.root.find(val, rbt.cmp)
	if !ok {
		return del, false
	}

	val = rbnDelete.Val
	rbt.Count--

	if rbt.Count == 0 {
		rbt.root = nil
		rbt.Min = nil
		rbt.Max = nil

		return val, true
	}

	if rbt.cmp(val, rbt.Min.Val) == 0 {
		rbt.Min, _ = rbt.Min.Next()
	}

	if rbt.cmp(val, rbt.Max.Val) == 0 {
		rbt.Max, _ = rbt.Max.Prev()
	}

	rbt.deleteCheckChildren(rbnDelete)

	return val, true
}

// deleteCheckChildren is the continuation of the Delete function (split for readability).
func (rbt *RBTree[T]) deleteCheckChildren(rbnDelete *RBNode[T]) {
	switch {
	case rbnDelete.left == nil && rbnDelete.right == nil: // no children
		rbt.deleteNoChildren(rbnDelete)
	case rbnDelete.left == nil: // one child
		rbnDelete.Val = rbnDelete.right.Val
		rbnDelete.right = nil
	case rbnDelete.right == nil:
		rbnDelete.Val = rbnDelete.left.Val
		rbnDelete.left = nil
	default: // left and right: find the next closest value, swap values, delete leaf
		rbnDelete.Val = rbt.findAndDeleteLeftmost(rbnDelete.right) // find and delete the leftmost successor of the right child
	}

	if rbt.cmp(rbnDelete.Val, rbt.Min.Val) == 0 {
		rbt.Min = rbnDelete
	}

	if rbt.cmp(rbnDelete.Val, rbt.Max.Val) == 0 {
		rbt.Max = rbnDelete
	}
}

// rotateRight moves the node down to the right.
//
//	    a             b
//	   / \           / \
//	  b   e   =>    c   a
//	 / \               / \
//	c   d             d   e
func (rbt *RBTree[T]) rotateRight(rbn *RBNode[T]) {
	if rbt.root == rbn {
		rbt.root = rbn.left
	}

	rbn.left.parent = rbn.parent
	rbn.parent = rbn.left
	rbn.left = rbn.parent.right

	if rbn.left != nil {
		rbn.left.parent = rbn
	}

	rbn.parent.right = rbn

	if rbn.parent.parent != nil { // update topnode parent child.
		if rbn.parent.parent.left == rbn {
			rbn.parent.parent.left = rbn.parent
		} else {
			rbn.parent.parent.right = rbn.parent
		}
	}
}

// rotateLeft moves the node down to the left.
//
//	  a             b
//	 / \           / \
//	e   b    =>   a   d
//	   / \       / \
//	  c   d     e   c
func (rbt *RBTree[T]) rotateLeft(rbn *RBNode[T]) {
	if rbt.root == rbn {
		rbt.root = rbn.right
	}

	rbn.right.parent = rbn.parent
	rbn.parent = rbn.right
	rbn.right = rbn.parent.left

	if rbn.right != nil {
		rbn.right.parent = rbn
	}

	rbn.parent.left = rbn

	if rbn.parent.parent != nil { // update topnode parent child.
		if rbn.parent.parent.left == rbn {
			rbn.parent.parent.left = rbn.parent
		} else {
			rbn.parent.parent.right = rbn.parent
		}
	}
}

// solveDoubleRed maintains the validity of the red-black tree if a red node has a red child.
func (rbt *RBTree[T]) solveDoubleRed(rbn *RBNode[T]) {
	switch {
	case isBlack(rbn.parent.left): // if sibling is left and black
		if !isBlack(rbn.left) { // making "line" from "right-triangle"
			rbt.rotateRight(rbn)
			rbn = rbn.parent
		}

		rbn.parent.isBlack = false
		rbn.isBlack = true

		rbt.rotateLeft(rbn.parent)
	case isBlack(rbn.parent.right): // if sibling is right and black
		if !isBlack(rbn.right) { // making "line" from "left-triangle"
			rbt.rotateLeft(rbn)
			rbn = rbn.parent
		}

		rbn.parent.isBlack = false
		rbn.isBlack = true

		rbt.rotateRight(rbn.parent)
	default: // if sibling is red
		rbn.parent.left.isBlack = true
		rbn.parent.right.isBlack = true

		if rbn.parent.parent != nil {
			rbn.parent.isBlack = false
			if !rbn.parent.parent.isBlack {
				rbt.solveDoubleRed(rbn.parent.parent)
			}
		}
	}
}

// solveDoubleBlack maintains the validity of the red-black tree after deletion.
func (rbt *RBTree[T]) solveDoubleBlack(rbn *RBNode[T]) {
	if rbt.root == rbn {
		return
	}

	parent := rbn.parent

	var (
		siblingIsRight bool
		sibling        *RBNode[T]
	)

	if parent.left == rbn || (parent.right != nil && parent.right != rbn) { // right sibling
		siblingIsRight = true
		sibling = parent.right
	} else { // left sibling
		sibling = parent.left
	}

	if sibling != nil && !sibling.isBlack { // red sibling
		parent.isBlack = false
		sibling.isBlack = true

		if siblingIsRight {
			rbt.rotateLeft(parent)
			sibling = parent.right
		} else {
			rbt.rotateRight(parent)
			sibling = parent.left
		}
	}

	// black sibling with black children
	if sibling.isBlack && isBlack(sibling.left) && isBlack(sibling.right) {
		sibling.isBlack = false

		if parent.isBlack {
			rbt.solveDoubleBlack(parent)

			return
		}

		parent.isBlack = true

		return
	}

	// black sibling with red child
	rbt.doubleBlackBlackSiblingRedChild(parent, sibling, siblingIsRight)
}

// doubleBlackBlackSiblingRedChild is the continuation of the solveDoubleBlack function (split for readability).
func (rbt *RBTree[T]) doubleBlackBlackSiblingRedChild(parent *RBNode[T], sibling *RBNode[T], siblingIsRight bool) {
	rightIsBlack := isBlack(sibling.right)
	leftIsBlack := isBlack(sibling.left)

	if rightIsBlack == siblingIsRight && leftIsBlack != siblingIsRight { // 🦊
		sibling.isBlack = false

		if siblingIsRight {
			sibling.left.isBlack = true
			rbt.rotateRight(sibling)
			sibling = parent.right
		} else {
			sibling.right.isBlack = true
			rbt.rotateLeft(sibling)
			sibling = parent.left
		}

		rightIsBlack = isBlack(sibling.right)
		leftIsBlack = isBlack(sibling.left)
	}

	sibling.isBlack = parent.isBlack
	parent.isBlack = true

	if siblingIsRight && !rightIsBlack {
		sibling.right.isBlack = true

		rbt.rotateLeft(parent)
	} else if !siblingIsRight && !leftIsBlack {
		sibling.left.isBlack = true

		rbt.rotateRight(parent)
	}
}

// findAndDeleteLeftmost deletes the leftmost node and returns its value.
func (rbt *RBTree[T]) findAndDeleteLeftmost(rbn *RBNode[T]) T {
	if rbn.left != nil {
		return rbt.findAndDeleteLeftmost(rbn.left)
	}

	if rbn.right != nil {
		rbn.right.parent = rbn.parent
		rbn.right.isBlack = true

		if rbn.parent.left == rbn {
			rbn.parent.left = rbn.right
		} else {
			rbn.parent.right = rbn.right
		}

		return rbn.Val
	}

	rbt.deleteNoChildren(rbn)

	return rbn.Val
}

// deleteNoChildren deletes a node without children.
func (rbt *RBTree[T]) deleteNoChildren(rbn *RBNode[T]) {
	if rbn.parent.left == rbn {
		rbn.parent.left = nil
	} else {
		rbn.parent.right = nil
	}

	if rbn.isBlack {
		rbt.solveDoubleBlack(rbn)
	}
}

package rbtree

import (
	"cmp"
	"math/rand/v2"
	"testing"
)

/*
		        70b
		       /   \
			50r     80r
		   / \      / \
	     20b  60b  75b 100b
*/
func initRBTBefore() *RBTree[int] {
	rbtBefore := &RBTree[int]{
		root: &RBNode[int]{
			Val:     70,
			isBlack: true,
		},
		cmp:   cmp.Compare[int],
		Count: 7,
	}

	rbtBefore.root.left = &RBNode[int]{
		Val:     50,
		isBlack: false,
		parent:  rbtBefore.root,
	}

	rbtBefore.root.left.left = &RBNode[int]{
		Val:     20,
		isBlack: true,
		parent:  rbtBefore.root.left,
	}

	rbtBefore.root.left.right = &RBNode[int]{
		Val:     60,
		isBlack: true,
		parent:  rbtBefore.root.left,
	}

	rbtBefore.root.right = &RBNode[int]{
		Val:     80,
		isBlack: false,
		parent:  rbtBefore.root,
	}

	rbtBefore.root.right.left = &RBNode[int]{
		Val:     75,
		isBlack: true,
		parent:  rbtBefore.root.right,
	}

	rbtBefore.root.right.right = &RBNode[int]{
		Val:     100,
		isBlack: true,
		parent:  rbtBefore.root.right,
	}

	rbtBefore.Min = rbtBefore.root.left.left
	rbtBefore.Max = rbtBefore.root.right.right

	return rbtBefore
}

func TestIsValid(t *testing.T) {
	t.Parallel()

	t.Run("IsValid: no cmp", func(t *testing.T) {
		t.Parallel()

		rbt := &RBTree[int]{}

		if rbt.IsValid() {
			t.Fail()
		}
	})

	t.Run("IsValid: empty tree", func(t *testing.T) {
		t.Parallel()

		rbt := &RBTree[int]{
			cmp: cmp.Compare[int],
		}

		if !rbt.IsValid() {
			t.Fail()
		}
	})

	t.Run("IsValid: red root", func(t *testing.T) {
		t.Parallel()

		rbt := &RBTree[int]{
			root: &RBNode[int]{
				Val:     70,
				isBlack: false,
			},
			cmp:   cmp.Compare[int],
			Count: 1,
		}

		if rbt.IsValid() {
			t.Fail()
		}
	})

	t.Run("IsValid: wrong count", func(t *testing.T) {
		t.Parallel()

		rbt := initRBTBefore()
		rbt.Count--

		if rbt.IsValid() {
			t.Fail()
		}
	})

	t.Run("IsValid: invalid black height", func(t *testing.T) {
		t.Parallel()

		rbt := initRBTBefore()

		rbt.root.left.isBlack = true

		if rbt.IsValid() {
			t.Fail()
		}
	})

	t.Run("IsValid: red child, red parent", func(t *testing.T) {
		t.Parallel()

		rbt := initRBTBefore()

		rbt.root.left.isBlack = true
		rbt.root.left.left.isBlack = false
		rbt.root.left.right.isBlack = false
		rbt.Count++

		rbt.root.left.left.left = &RBNode[int]{
			Val:     10,
			isBlack: false,
			parent:  rbt.root.left.left,
		}

		if rbt.IsValid() {
			t.Fail()
		}
	})

	t.Run("IsValid: wrong parent, left child", func(t *testing.T) {
		t.Parallel()

		rbt := initRBTBefore()

		rbt.root.left.parent = rbt.root.right

		if rbt.IsValid() {
			t.Fail()
		}
	})

	t.Run("IsValid: wrong parent, right child", func(t *testing.T) {
		t.Parallel()

		rbt := initRBTBefore()

		rbt.root.right.parent = rbt.root.left

		if rbt.IsValid() {
			t.Fail()
		}
	})
}

func TestClone(t *testing.T) {
	t.Parallel()

	t.Run("Clone: empty tree", func(t *testing.T) {
		t.Parallel()

		rbt := &RBTree[int]{
			cmp: cmp.Compare[int],
		}

		rbtCloned := rbt.Clone()

		if rbtCloned.root != nil || rbtCloned.Count != 0 || rbtCloned.cmp == nil || rbtCloned.Max != nil || rbtCloned.Min != nil {
			t.Fail()
		}
	})

	t.Run("Clone: 3-node tree", func(t *testing.T) {
		t.Parallel()

		rbt := &RBTree[int]{
			root: &RBNode[int]{
				Val:     70,
				isBlack: true,
			},
			cmp:   cmp.Compare[int],
			Count: 3,
		}

		rbt.root.left = &RBNode[int]{
			Val:     50,
			isBlack: true,
			parent:  rbt.root,
		}

		rbt.root.right = &RBNode[int]{
			Val:     80,
			isBlack: true,
			parent:  rbt.root,
		}

		rbt.Max = rbt.root.right
		rbt.Min = rbt.root.left

		rbtCloned := rbt.Clone()

		if rbtCloned.root == nil || rbtCloned.Count != rbt.Count {
			t.Fail()
		}

		if rbtCloned.cmp(rbtCloned.Min.Val, rbt.Min.Val) != 0 || rbtCloned.cmp(rbtCloned.Max.Val, rbt.Max.Val) != 0 {
			t.Fail()
		}

		if rbtCloned.root.left == nil || rbtCloned.root.left.parent != rbtCloned.root {
			t.Fail()
		}

		if rbtCloned.root.right == nil || rbtCloned.root.right.parent != rbtCloned.root {
			t.Fail()
		}

		if rbtCloned.cmp(rbtCloned.root.Val, rbt.root.Val) != 0 {
			t.Fail()
		}

		if rbtCloned.cmp(rbtCloned.root.left.Val, rbt.root.left.Val) != 0 {
			t.Fail()
		}

		if rbtCloned.cmp(rbtCloned.root.right.Val, rbt.root.right.Val) != 0 {
			t.Fail()
		}

		if rbtCloned.root.isBlack != rbt.root.isBlack {
			t.Fail()
		}

		if rbtCloned.root.left.isBlack != rbt.root.left.isBlack {
			t.Fail()
		}

		if rbtCloned.root.right.isBlack != rbt.root.right.isBlack {
			t.Fail()
		}
	})
}

func TestEqualTo(t *testing.T) {
	t.Parallel()

	t.Run("EqualTo: empty trees", func(t *testing.T) {
		t.Parallel()

		rbt := &RBTree[int]{
			cmp: cmp.Compare[int],
		}
		anotherRBT := &RBTree[int]{
			cmp: cmp.Compare[int],
		}

		if !rbt.EqualTo(anotherRBT) {
			t.Fail()
		}
	})

	t.Run("EqualTo: empty and non-empty tree", func(t *testing.T) {
		t.Parallel()

		rbt := &RBTree[int]{
			cmp: cmp.Compare[int],
		}
		anotherRBT := initRBTBefore()

		if rbt.EqualTo(anotherRBT) {
			t.Fail()
		}
	})

	t.Run("EqualTo: nil and non-nil tree", func(t *testing.T) {
		t.Parallel()

		rbt := &RBTree[int]{
			cmp: cmp.Compare[int],
		}

		var anotherRBT *RBTree[int]

		if rbt.EqualTo(anotherRBT) {
			t.Fail()
		}
	})

	t.Run("EqualTo: non-equal structure", func(t *testing.T) {
		t.Parallel()

		rbt := initRBTBefore()
		anotherRBT := initRBTBefore()

		rbt.root.left.left.right = &RBNode[int]{
			Val:     10,
			isBlack: false,
			parent:  rbt.root.left.left,
		}

		anotherRBT.root.right.right.right = &RBNode[int]{
			Val:     110,
			isBlack: false,
			parent:  rbt.root.right.right,
		}

		rbt.Count++
		anotherRBT.Count++

		if rbt.EqualTo(anotherRBT) {
			t.Fail()
		}
	})

	t.Run("EqualTo: non-equal structure - left node", func(t *testing.T) {
		t.Parallel()

		rbt := initRBTBefore()
		anotherRBT := initRBTBefore()

		anotherRBT.root.left.left.left = &RBNode[int]{
			Val:     10,
			isBlack: false,
			parent:  anotherRBT.root.left.left,
		}

		rbt.root.right.right.right = &RBNode[int]{
			Val:     110,
			isBlack: false,
			parent:  rbt.root.right.right,
		}

		rbt.Count++
		anotherRBT.Count++

		if rbt.EqualTo(anotherRBT) {
			t.Fail()
		}
	})

	t.Run("EqualTo: non-equal structure - right node", func(t *testing.T) {
		t.Parallel()

		rbt := initRBTBefore()
		anotherRBT := initRBTBefore()

		anotherRBT.root.left.left.right = &RBNode[int]{
			Val:     10,
			isBlack: false,
			parent:  anotherRBT.root.left.left,
		}

		rbt.root.right.right.right = &RBNode[int]{
			Val:     110,
			isBlack: false,
			parent:  rbt.root.right.right,
		}

		rbt.Count++
		anotherRBT.Count++

		if rbt.EqualTo(anotherRBT) {
			t.Fail()
		}
	})

	t.Run("EqualTo: different size", func(t *testing.T) {
		t.Parallel()

		rbt := initRBTBefore()
		anotherRBT := initRBTBefore()

		rbt.root.left.left.right = &RBNode[int]{
			Val:     10,
			isBlack: false,
			parent:  rbt.root.left.left,
		}

		rbt.Count++

		if rbt.EqualTo(anotherRBT) {
			t.Fail()
		}
	})

	t.Run("EqualTo: different values", func(t *testing.T) {
		t.Parallel()

		rbt := initRBTBefore()
		anotherRBT := initRBTBefore()

		anotherRBT.root.Val = 75

		if rbt.EqualTo(anotherRBT) {
			t.Fail()
		}
	})

	t.Run("EqualTo: equal trees", func(t *testing.T) {
		t.Parallel()

		rbt := initRBTBefore()
		anotherRBT := initRBTBefore()

		if !rbt.EqualTo(anotherRBT) {
			t.Fail()
		}
	})
}

func TestString(t *testing.T) {
	t.Parallel()

	t.Run("String: empty tree", func(t *testing.T) {
		t.Parallel()

		rbt := &RBTree[int]{
			cmp: cmp.Compare[int],
		}

		if rbt.String() != "" {
			t.Fail()
		}
	})

	t.Run("String: non-empty tree", func(t *testing.T) {
		t.Parallel()

		rbt := initRBTBefore()
		expectedResult := "   100\n  80\n   75\n 70\n   60\n  50\n   20\n"

		if rbt.String() != expectedResult {
			t.Fail()
		}
	})
}

func TestNext(t *testing.T) {
	t.Parallel()

	t.Run("Next: root-only tree", func(t *testing.T) {
		t.Parallel()

		rbt := &RBTree[int]{
			root: &RBNode[int]{
				Val:     20,
				isBlack: true,
			},
		}

		node, ok := rbt.root.Next()
		if ok || node != nil {
			t.Fail()
		}
	})

	t.Run("Next: is a right child, has a bigger grandparent", func(t *testing.T) {
		t.Parallel()

		rbt := initRBTBefore()

		node, ok := rbt.root.left.right.Next()
		if !ok || node != rbt.root {
			t.Fail()
		}
	})

	t.Run("Next: is a right child", func(t *testing.T) {
		t.Parallel()

		rbt := initRBTBefore()

		node, ok := rbt.root.right.right.Next()
		if ok || node != nil {
			t.Fail()
		}
	})
}

func TestPrev(t *testing.T) {
	t.Parallel()

	t.Run("Prev: root-only tree", func(t *testing.T) {
		t.Parallel()

		rbt := &RBTree[int]{
			root: &RBNode[int]{
				Val:     20,
				isBlack: true,
			},
		}

		node, ok := rbt.root.Prev()
		if ok || node != nil {
			t.Fail()
		}
	})

	t.Run("Prev: is a left child, has a smaller grandparent", func(t *testing.T) {
		t.Parallel()

		rbt := initRBTBefore()

		node, ok := rbt.root.right.left.Prev()
		if !ok || node != rbt.root {
			t.Fail()
		}
	})

	t.Run("Prev: is a left child", func(t *testing.T) {
		t.Parallel()

		rbt := initRBTBefore()

		node, ok := rbt.root.left.left.Prev()
		if ok || node != nil {
			t.Fail()
		}
	})
}

func TestFind(t *testing.T) {
	t.Parallel()

	t.Run("Find: empty tree", func(t *testing.T) {
		t.Parallel()

		node, ok := (&RBTree[int]{}).Find(10)
		if ok || node != nil {
			t.Fail()
		}
	})

	t.Run("Find: non-existent value", func(t *testing.T) {
		t.Parallel()

		rbt := &RBTree[int]{
			root: &RBNode[int]{
				Val:     20,
				isBlack: true,
			},
			cmp: cmp.Compare[int],
		}

		node, ok := rbt.Find(10)
		if ok || node != nil {
			t.Fail()
		}
	})
}

func TestInsert(t *testing.T) {
	t.Parallel()

	t.Run("Insert: existent value", func(t *testing.T) {
		t.Parallel()

		const val = 75

		rbt := initRBTBefore()

		existentNode, ok := rbt.Insert(val)
		if ok || !rbt.IsValid() || existentNode == nil || rbt.cmp(existentNode.Val, val) != 0 {
			t.Fail()
		}
	})
}

func TestDelete(t *testing.T) {
	t.Parallel()

	t.Run("Delete: the last node", func(t *testing.T) {
		t.Parallel()

		const val = 50

		rbtBefore := &RBTree[int]{
			root: &RBNode[int]{
				Val:     val,
				isBlack: true,
			},
			cmp:   cmp.Compare[int],
			Count: 1,
		}

		rbtBefore.Max = rbtBefore.root
		rbtBefore.Min = rbtBefore.root

		_, ok := rbtBefore.Delete(val)
		if !ok || !rbtBefore.IsValid() {
			t.Fail()
		}
	})

	t.Run("Delete: empty tree", func(t *testing.T) {
		t.Parallel()

		rbtBefore := &RBTree[int]{
			cmp: cmp.Compare[int],
		}

		_, ok := rbtBefore.Delete(10)
		if ok || !rbtBefore.IsValid() {
			t.Fail()
		}
	})

	t.Run("Delete: non-existent value", func(t *testing.T) {
		t.Parallel()

		rbtBefore := &RBTree[int]{
			root: &RBNode[int]{
				Val:     20,
				isBlack: true,
			},
			cmp:   cmp.Compare[int],
			Count: 1,
		}

		rbtBefore.Min = rbtBefore.root
		rbtBefore.Max = rbtBefore.root

		_, ok := rbtBefore.Delete(30)
		if ok || !rbtBefore.IsValid() {
			t.Fail()
		}
	})
}

func TestRandomInsertDelete(t *testing.T) {
	t.Parallel()

	const (
		maxTreeSize = 1000
		iterations  = 1000
	)

	for range iterations {
		rbt := NewOrdered[int]()
		treeSize := rand.IntN(maxTreeSize) + 1
		insertedValues := make(map[int]struct{})

		for range treeSize {
			inserted, ok := rbt.Insert(rand.Int())
			if !ok {
				continue
			}

			insertedValues[inserted.Val] = struct{}{}

			if rbt.Count != len(insertedValues) || !rbt.IsValid() {
				t.FailNow()
			}
		}

		numItemsToDelete := rand.IntN(len(insertedValues))

		for range numItemsToDelete {
			var randVal int

			for k := range insertedValues {
				randVal = k
				delete(insertedValues, k)

				break
			}

			_, _ = rbt.Delete(randVal)

			if rbt.Count != len(insertedValues) || !rbt.IsValid() {
				t.FailNow()
			}
		}
	}
}

func BenchmarkRW(b *testing.B) {
	treeSizes := map[string]int{
		"1000":     1000,
		"100000":   100000,
		"10000000": 10000000,
	}

	for name, treeSize := range treeSizes {
		rbt := NewOrdered[int]()

		for i := range treeSize {
			_, _ = rbt.Insert(i)
		}

		b.Run("Find-"+name, func(b *testing.B) {
			for range b.N {
				_, _ = rbt.Find(treeSize)
			}
		})
	}

	for name, treeSize := range treeSizes {
		rbt := NewOrdered[int]()

		b.Run("InsertDelete-"+name, func(b *testing.B) {
			for range b.N {
				for i := range treeSize {
					_, _ = rbt.Insert(i)
				}

				for i := range treeSize {
					_, _ = rbt.Delete(i)
				}
			}
		})
	}
}

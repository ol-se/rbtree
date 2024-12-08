package main

import (
	"cmp"
	"fmt"
	"math/rand/v2"

	"github.com/ol-se/rbtree"
)

type orderedMapItem[K cmp.Ordered, V any] struct {
	key K
	val V
}

func compareKeys[K cmp.Ordered, V any](first orderedMapItem[K, V], second orderedMapItem[K, V]) int {
	return cmp.Compare(first.key, second.key)
}

func main() {
	rbTree := rbtree.New(compareKeys[int, int])

	for range 10 {
		rbNode, ok := rbTree.Insert(orderedMapItem[int, int]{
			key: rand.IntN(10),
			val: rand.IntN(1000),
		})

		if ok {
			fmt.Printf("Item %v inserted.\n", rbNode.Val)
		} else {
			fmt.Printf("Key %d already exists with value %d.\n", rbNode.Val.key, rbNode.Val.val)
		}
	}

	const key = 5

	fmt.Printf("Searching for a specific key: %d.\n", key)

	rbNode, ok := rbTree.Find(orderedMapItem[int, int]{
		key: key,
	})

	if ok {
		fmt.Printf("Key %d found with value %d.\n", rbNode.Val.key, rbNode.Val.val)
	} else {
		fmt.Printf("Key %d not found.\n", key)
	}

	fmt.Println("Traversing map from Max to Min:")

	for rbNode, ok := rbTree.Max, rbTree.Max != nil; ok; rbNode, ok = rbNode.Prev() {
		fmt.Printf("%v ", rbNode.Val)
	}

	fmt.Println("\nCloning map.")

	rbTreeCopy := rbTree.Clone()

	for range 5 {
		randKey := rand.IntN(10)

		item, ok := rbTreeCopy.Delete(orderedMapItem[int, int]{
			key: randKey,
		})

		if ok {
			fmt.Printf("Item %v deleted from the clone.\n", item)
		} else {
			fmt.Printf("Key %d not found in the clone.\n", randKey)
		}
	}

	if rbTreeCopy.IsValid() {
		fmt.Println("The clone is valid.")

		if !rbTreeCopy.EqualTo(rbTree) {
			fmt.Println("The maps are not equal.")
		} else {
			fmt.Println("The maps are equal.")
		}
	} else {
		fmt.Println("The clone is invalid!")
	}
}

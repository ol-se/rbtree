package main

import (
	"fmt"
	"math/rand/v2"

	"github.com/ol-se/rbtree"
)

func main() {
	rbTree := rbtree.NewOrdered[int]()

	for range 10 {
		rbNode, ok := rbTree.Insert(rand.IntN(10))

		if ok {
			fmt.Printf("Item %d inserted.\n", rbNode.Val)
		} else {
			fmt.Printf("Key %d already exists.\n", rbNode.Val)
		}
	}

	fmt.Println("Traversing tree from Min to Max:")

	for rbNode, ok := rbTree.Min, rbTree.Min != nil; ok; rbNode, ok = rbNode.Next() {
		fmt.Printf("%d ", rbNode.Val)
	}

	fmt.Printf("Size: %d\nTree:\n%s", rbTree.Count, rbTree)
}

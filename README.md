# Golang Generic Red-Black Tree

Package rbtree is a zero-dependencies library that provides methods to work with generic [red-black tree](https://en.wikipedia.org/wiki/Red%E2%80%93black_tree). Both primitives and user-defined types can be used as values of the red-black tree nodes.

## Go version
1.22+

## Usage

```go
package main

import (
	"fmt"

	"github.com/ol-se/rbtree"
)

func main() {
	t := rbtree.NewOrdered[int]()

	for i := range 10 {
		_, _ = t.Insert(i)
	}

	fmt.Println(t)
}
```

For more examples on how to use (e.g. iterate the tree or work with user-defined types), see examples.

## Complexity

- Insert: `O(log n)`
- Delete: `O(log n)`

```
BenchmarkRW/InsertDelete-1000-4       2908       626676 ns/op
BenchmarkRW/InsertDelete-100000-4       14     74021303 ns/op
BenchmarkRW/InsertDelete-10000000-4      1   7736143898 ns/op
```

- Find: `O(log n)`

```
BenchmarkRW/Find-1000-4       8200576   150.4 ns/op
BenchmarkRW/Find-100000-4     5125063   234.4 ns/op
BenchmarkRW/Find-10000000-4   3669334   328.7 ns/op
```

## License

Distributed under the [Unlicense](https://unlicense.org/).
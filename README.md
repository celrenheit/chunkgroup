# Chunkgroup

Package chunkgroup provides a way to schedule the execution of a function when a chunk of data is ready.

## Example

```go
package main

import (
	"fmt"
	"log"

	"github.com/celrenheit/chunkgroup"
)

func main() {
	cg := chunkgroup.New(10, 5, func(items []int) error {
		fmt.Println(items)
		return nil
	})

	for i := range 105 {
		cg.Add(i)
	}

	if err := cg.Flush(); err != nil {
		log.Fatal(err)
	}
}
```

This should output something like this (order may be different due to concurrency being over 1):
```
[40 41 42 43 44 45 46 47 48 49]
[30 31 32 33 34 35 36 37 38 39]
[50 51 52 53 54 55 56 57 58 59]
[10 11 12 13 14 15 16 17 18 19]
[70 71 72 73 74 75 76 77 78 79]
[0 1 2 3 4 5 6 7 8 9]
[60 61 62 63 64 65 66 67 68 69]
[20 21 22 23 24 25 26 27 28 29]
[100 101 102 103 104]
[90 91 92 93 94 95 96 97 98 99]
[80 81 82 83 84 85 86 87 88 89]
```
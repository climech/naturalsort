# naturalsort #

Natural sorting implemented in Go.

The package expects valid UTF-8-encoded strings. If you're accepting user input, check with [`utf8.Valid`](https://golang.org/pkg/unicode/utf8/#Valid) or [`utf8.ValidString`](https://golang.org/pkg/unicode/utf8/#ValidString) first.  The strings may contain arbitrarily large numbers.

## Usage ##

A simple slice of strings:

```go
package main

import (
    "fmt"

    "github.com/climech/naturalsort"
)

func main() {
  items := []string{"Chapter 1", "Chapter 11", "Chapter 2"}
  naturalsort.Sort(items)
  
  for _, item := range items {
    fmt.Println(item)
  }
}
```

Output:

```
Chapter 1
Chapter 2
Chapter 11
```

A slice of structs sorted by field:

```go
package main

import (
	"fmt"
	"sort"

	"github.com/climech/naturalsort"
)

type Movie struct {
	ID    uint
	Title string
}

func main() {
	movies := []*Movie{
		{24, "Die Hard 2"},
		{17, "Die Hard 1"},
		{42, "Die Hard 10"},
	}

	sort.SliceStable(movies, func(i, j int) bool {
		return naturalsort.Compare(movies[i].Title, movies[j].Title)
	})

	for _, movie := range movies {
		fmt.Printf("%s (ID: %d)\n", movie.Title, movie.ID)
	}
}
```

Output:

```
Die Hard 1 (ID: 17)
Die Hard 2 (ID: 24)
Die Hard 10 (ID: 42)
```


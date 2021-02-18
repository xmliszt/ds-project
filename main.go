package main

import (
	"fmt"

	"github.com/xmliszt/e-safe/foo"
)

func main() {
	fmt.Println(foo.Foo(), foo.Sum(1,2))
}
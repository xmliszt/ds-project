package foo

import (
	"fmt"

	esafe "github.com/xmliszt/e-safe"
)

func Foo() string {
	fmt.Println(esafe.Sum(1,2))
	return "hello"
}
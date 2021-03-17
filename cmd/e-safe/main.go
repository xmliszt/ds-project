package main

import (
	"fmt"
	"os"

	"github.com/xmliszt/e-safe/pkg/locksmith"
)

func main() {
	err := locksmith.Start()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

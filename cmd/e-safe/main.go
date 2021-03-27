package main

import (
	"log"
	"os"

	"github.com/xmliszt/e-safe/pkg/locksmith"
)

func main() {

	err := locksmith.Start()
	if err != nil {
		log.Fatalln(err)
		os.Exit(1)
	}
}

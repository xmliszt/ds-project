package main

import (
	"github.com/xmliszt/e-safe/pkg/locksmith"
)

func main() {
	locksmith.InitializeLocksmith(5)
}
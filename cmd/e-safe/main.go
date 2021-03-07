package main

import (
	"github.com/xmliszt/e-safe/pkg/locksmith"
)

func main() {
	locksmithServer := locksmith.InitializeLocksmith(5)
	go locksmithServer.Node.HandleMessageReceived()
	locksmithServer.StartAllNodes()
	go locksmithServer.MonitorNodes()
}
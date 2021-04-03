package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/xmliszt/e-safe/pkg/file"
	"github.com/xmliszt/e-safe/pkg/locksmith"
	"github.com/xmliszt/e-safe/pkg/node"
	"github.com/xmliszt/e-safe/pkg/register"
)

func main() {

	if err := file.CreateStoragePath(); err != nil {
		log.Fatal(err)
	}
	register.Regsiter()

	var role string
	flag.StringVar(&role, "r", "node", "Select role of the process: node, locksmith")
	flag.Parse()

	if role != "node" && role != "locksmith" {
		log.Fatalln("Only support role: [node, locksmith]!")
	}

	switch role {
	case "node":
		var nodeID int
		fmt.Print("Enter Node ID to start (>=1): ")
		fmt.Scan(&nodeID)
		if nodeID < 1 {
			log.Fatalln("Node number must be larger than 0!")
		}
		log.Printf("Node %d start!\n", nodeID)
		node.Start(nodeID)
	case "locksmith":
		log.Println("Locksmith start!")
		locksmith.Start()
	}
}

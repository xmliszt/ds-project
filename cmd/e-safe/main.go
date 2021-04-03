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

	var isNode bool
	var isLocksmith bool
	flag.BoolVar(&isNode, "node", false, "Start as a node")
	flag.BoolVar(&isLocksmith, "locksmith", false, "Start as locksmith")
	flag.Parse()

	if isNode {
		var nodeID int
		fmt.Print("Enter Node ID to start (>=1): ")
		fmt.Scan(&nodeID)
		if nodeID < 1 {
			log.Fatalln("Node number must be larger than 0!")
		}
		err := file.CreateNodeStoragePath(nodeID)
		if err != nil {
			log.Fatal(err)
		}
		log.Printf("Node %d start!\n", nodeID)
		node.Start(nodeID)
	}

	if isLocksmith {
		log.Println("Locksmith start!")
		locksmith.Start()
	}

	fmt.Println("Please select a mode to start: -node, -locksmith")
}

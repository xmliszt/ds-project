package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/xmliszt/e-safe/config"
	"github.com/xmliszt/e-safe/pkg/locksmith"
	"github.com/xmliszt/e-safe/pkg/node"
	"github.com/xmliszt/e-safe/pkg/register"
)

func main() {

	register.Regsiter()

	var role string
	flag.StringVar(&role, "r", "node", "Select role of the process: node, locksmith")
	flag.Parse()

	if role != "node" && role != "locksmith" {
		log.Fatalln("Only support role: [node, locksmith]!")
	}

	config, err := config.GetConfig()
	if err != nil {
		log.Fatal(err)
	}

	switch role {
	case "node":
		var nodeID int
		fmt.Printf("Enter Node ID to start (1 ~ %d): ", config.ConfigNode.Number)
		fmt.Scan(&nodeID)
		if nodeID < 1 || nodeID > config.ConfigNode.Number {
			log.Fatalln("If you want to add more nodes, please change the node number in config.yaml!")
		}
		log.Printf("Node %d start!\n", nodeID)
		node.Start(nodeID)
	case "locksmith":
		log.Println("Locksmith start!")
		locksmith.Start()
	}
}

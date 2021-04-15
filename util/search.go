package util

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

func MapHashToVNode(virtualNodeMap map[int]string, virtualNodeLocation []int, hashedValue uint32) string {
	var virtualNodeName string
	for location := range virtualNodeLocation {
		if int(hashedValue) < location {
			virtualNodeName = virtualNodeMap[location]
			break
		}
		continue
	}
	return virtualNodeName
}

func MapHashToVNodeLoc(virtualNodeMap map[int]string, virtualNodeLocation []int, hashedValue uint32) int {
	var result int
	fmt.Println("These are the inputs", virtualNodeLocation, virtualNodeMap, hashedValue)
	fmt.Println("this is the int version of the hash", int(hashedValue))
	intHashedValue := int(hashedValue)
	for idx, location := range virtualNodeLocation {
		if intHashedValue < location {
			fmt.Println("This is the location in a loop", location)
			result = virtualNodeLocation[idx]
			break
		}
		continue
	}
	return result
}

func NodePidFromVNode(virtualNodeName string) int {

	splitVNodeName := strings.Split(virtualNodeName, "-")
	pid, err := strconv.Atoi(splitVNodeName[0])
	if err != nil {
		log.Println(err)
	}

	return pid
}

func FindNextVNode(ringStruct []int, virtualNodeMap map[int]string, virtualNodeLocation []int, hashedValue uint32) string {
	var nextVirtualNode string

	for idx, location := range virtualNodeLocation {
		if int(hashedValue) < location {
			nextVirtualNode = virtualNodeMap[ringStruct[(idx+1)]]
			break
		}
		continue
	}
	return nextVirtualNode
}

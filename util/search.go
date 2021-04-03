package util

import (
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

func NodePidFromVNode(virtualNodeName string) int {

	splitVNodeName := strings.Split(virtualNodeName, "-")
	pid, err := strconv.Atoi(splitVNodeName[0])

	if err != nil {
		log.Println("Error: get PID from VNode")
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

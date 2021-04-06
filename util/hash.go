package util

import (
	"hash/fnv"
	"log"
	"strconv"
)

func GetHash(s string) (uint32, error) {
	h := fnv.New32a()
	_, err := h.Write([]byte(s))
	if err != nil {
		return 0, err
	}
	return h.Sum32(), nil
}
func StringHashToUint(hashedValue string) uint32 {
	iHashedValue, intConvertError := strconv.Atoi(hashedValue)
	if intConvertError != nil {
		log.Println("Error: Converting hash from string to uint32")
	}
	uHashedValue := uint32(iHashedValue)
	return uHashedValue
}

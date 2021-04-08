package util

import (
	"hash/fnv"
	"log"
	"strconv"

	"github.com/xmliszt/e-safe/config"
)

func GetHash(s string) (uint32, error) {
	config, err := config.GetConfig()
	if err != nil {
		return 0, err
	}

	var un uint32
	for i := 0; i < config.NumberOfHashing; i++ {
		h := fnv.New32a()
		_, err := h.Write([]byte(s))
		if err != nil {
			return 0, err
		}
		un = h.Sum32()
		s = strconv.Itoa(int(un))
	}
	return un, nil
}

func StringHashToUint(hashedValue string) uint32 {
	iHashedValue, intConvertError := strconv.Atoi(hashedValue)
	if intConvertError != nil {
		log.Println("Error: Converting hash from string to uint32")
	}
	uHashedValue := uint32(iHashedValue)
	return uHashedValue
}

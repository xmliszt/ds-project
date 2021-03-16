package util

import "fmt"

// SliceRemoveInt removes an integer value from a slice of integer, if value does not exist, it will throw an error
func SliceRemoveInt(s []int, v int) ([]int, error) {
	for idx, i := range s {
		if i == v {
			r := append(s[:idx], s[idx+1:]...)
			return r, nil
		}
	}
	return nil, fmt.Errorf("Slice %v does not contain %v", s, v)
}


// IntInSlice check if a given integer v is in the slice of integers s
func IntInSlice(s []int, v int) bool {
	for _, i := range s {
		if i == v {
			return true
		}
	}
	return false
}

// Find max value in list
func FindMax(array []int) (int) {
    var max int = array[0]
    for _, value := range array {
        if max < value {
            max = value
        }
    }
    return max
}
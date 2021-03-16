package util

import (
	"reflect"
	"testing"
)

func TestSliceRemoveInt(t *testing.T) {
	a := []int{1,2,3,4,5}
	r, err := SliceRemoveInt(a, 3)	// Remove existent 3
	if err != nil {
		t.Error(err)
	}
	expected := []int{1,2,4,5}
	if !reflect.DeepEqual(r, expected) {
		t.Errorf("Resultant slice %v is not the same as expected %v", r, expected)
	}

	f, err := SliceRemoveInt(a, 6)	// Remove non-existent 6
	if err == nil {
		t.Errorf("Expect to throw error as 6 does not exist in slice, but instead got result %v", f)
	}
}

func TestIntInSlice(t *testing.T) {
	s := []int{1,2,3}
	a := IntInSlice(s, 1)
	b := IntInSlice(s, 4)
	if !a || b {
		t.Errorf("Check if 1 and 4 in slice [1 2 3] failed, results: (1)%v (4)%v", a, b)
	}
}

func TestFindMax(t *testing.T) {
	testArray := []int{1,2,3,5,4}
	max := FindMax(testArray)
	if max != 5 {
		t.Errorf("5 should be the max but it is : %d!", max)
	}
}
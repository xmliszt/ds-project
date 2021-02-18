package main

import "testing"

func TestSum(t *testing.T) {
	r := Sum(1, 2)
	if r != 1 + 2 {
		t.Log("1 + 2 should be 3 but got", r)
		t.Fail()
	}
}
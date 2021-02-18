package main

import (
	"testing"

	"github.com/xmliszt/e-safe/foo"
)

func TestSum(t *testing.T) {
	r := foo.Sum(1, 2)
	if r != 1 + 2 {
		t.Log("1 + 2 should be 3 but got", r)
		t.Fail()
	}
}
package foo

import "testing"

func TestFoo(t *testing.T) {
	r := Foo()
	if r != "hello" {
		t.Log("Foo should return hello but instead got", r)
	}
	t.Error("Test failed!")
}
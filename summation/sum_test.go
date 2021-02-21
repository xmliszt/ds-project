package summation

import "testing"

func TestSumTwo(t *testing.T) {
	a := 1
	b := 2
	r := SumTwo(a, b)
	if r != a + b {
		t.Errorf("%d + %d should be %d but instead got %d\n", a, b, a+b, r)
	}
}
package main

import "testing"

func sum(ops ...int) int {
	s := 0
	for _, op := range ops {
		s += op
	}
	return s
}

func TestParamsTransform(t *testing.T) {
	arr3 := [...]int{1, 3, 4, 5}
	for i := 0; i < len(arr3); i++ {
		t.Log(arr3[i])
	}
	t.Log(sum(1, 2, 3))
}

func TestDefer(t *testing.T) {
	defer func() {
		t.Log("Clear resources")
	}()
	t.Log("Started")
	//defer仍会执⾏行行
	panic("Fatal error")
}

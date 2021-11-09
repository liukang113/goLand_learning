package main

import "testing"

func TestSlice(t *testing.T) {
	var s0 []int
	s1 := append(s0, 1)
	t.Log(len(s1), cap(s1))

	s2 := []int{1, 2, 3, 4}
	t.Log(len(s2), cap(s2))

	// len 个元素被初始化为默认值，未初始化值不可访问
	s3 := make([]int, 3, 5)
	t.Log(len(s3), cap(s3))
}

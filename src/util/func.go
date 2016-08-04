package util

import (
	"math/rand"
	"sort"
)

func ContainsInt(l []int, v int) bool {
	for _, v1 := range l {
		if v1 == v {
			return true
		}
	}
	return false
}

func ContainsStr(l []string, v string) bool {
	for _, v1 := range l {
		if v1 == v {
			return true
		}
	}
	return false
}

type IntRandSlice []int

func (p IntRandSlice) Len() int           { return len(p) }
func (p IntRandSlice) Less(i, j int) bool { return rand.Float32() < 0.5 }
func (p IntRandSlice) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }

// Sort is a convenience method.
func (p IntRandSlice) Sort() { sort.Sort(p) }

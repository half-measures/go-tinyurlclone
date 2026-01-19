package main

import (
	"testing"
)

func TestGenerateSlug(t *testing.T) {
	lengths := []int{4, 6, 8, 10}
	for _, l := range lengths {
		slug := generateSlug(l)
		if len(slug) != l {
			t.Errorf("Expected slug length %d, got %d", l, len(slug))
		}
	}
}

// Copyright 2025 David Fritz. All rights reserved.
// This software may be modified and distributed under the terms of the BSD
// 2-clause license. See the LICENSE file for details.

package vdb

import "testing"

func TestNewVector(t *testing.T) {
	_, err := NewVector([]float64{1.5, 2.0}, "test")
	if err != nil {
		t.Fatal(err)
	}
}

func TestNorm(t *testing.T) {
	x, err := NewVector([]float64{1.5, 2.0, 1.234}, "test")
	if err != nil {
		t.Fatal(err)
	}
	if x.N != 2.7879662838707358 {
		t.Fatalf("invalid norm: %v", x.N)
	}
}

func TestSimilarity(t *testing.T) {
	a, err := NewVector([]float64{10, 9}, "kris is a potato")
	if err != nil {
		t.Fatal(err)
	}
	b, err := NewVector([]float64{5, 9}, "Who is a potato?")
	if err != nil {
		t.Fatal(err)
	}
	s, err := a.CosineSimilarity(b)
	if err != nil {
		t.Fatal(err)
	}
	if s != 0.9457559355278488 {
		t.Fatalf("incorrect similarity: %v. Expected 1.0", s)
	}
}

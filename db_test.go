// Copyright 2025 David Fritz. All rights reserved.
// This software may be modified and distributed under the terms of the BSD
// 2-clause license. See the LICENSE file for details.

package vdb

import (
	"fmt"
	"testing"
)

func TestDBSimilarity(t *testing.T) {
	a, err := NewVector([]float64{10, 9}, "kris is a potato")
	if err != nil {
		t.Fatal(err)
	}
	b, err := NewVector([]float64{-5, 0}, "john is a fool")
	if err != nil {
		t.Fatal(err)
	}
	c, err := NewVector([]float64{5000, 10000}, "fritz is a german")
	if err != nil {
		t.Fatal(err)
	}

	db := new(DB)
	db.AddVector(a)
	db.AddVector(b)
	db.AddVector(c)
	if db.Len() != 3 {
		t.Fatal("invalid length")
	}

	x, err := NewVector([]float64{8, 9}, "Who is a potato?")
	if err != nil {
		t.Fatal(err)
	}

	r, err := db.SimilarVectors(x, 1, 0.1)
	if err != nil {
		t.Fatal(err)
	}
	if r[0].GetData().(string) != "kris is a potato" {
		t.Fatalf("invalid vector returned: %v", r[0].GetData().(string))
	}
}

func TestDBSimilarity2(t *testing.T) {
	a, err := NewVector([]float64{10, 9}, "kris is a potato")
	if err != nil {
		t.Fatal(err)
	}
	b, err := NewVector([]float64{5, 1}, "john is a fool")
	if err != nil {
		t.Fatal(err)
	}
	c, err := NewVector([]float64{-5000, 10000}, "fritz is a german")
	if err != nil {
		t.Fatal(err)
	}
	d, err := NewVector([]float64{-5000, 12000}, "more data")
	if err != nil {
		t.Fatal(err)
	}

	db := new(DB)
	db.AddVector(b)
	db.AddVector(d)
	db.AddVector(c)
	db.AddVector(a)
	if db.Len() != 4 {
		t.Fatal("invalid length")
	}

	x, err := NewVector([]float64{8, 9}, "Who is a potato?")
	if err != nil {
		t.Fatal(err)
	}

	r, err := db.SimilarVectors(x, 4, 0.1)
	if err != nil {
		t.Fatal(err)
	}
	if r[0].GetData().(string) != "kris is a potato" {
		t.Fatalf("invalid vector returned: %v", r[0].GetData().(string))
	}
	if r[1].GetData().(string) != "john is a fool" {
		t.Fatalf("invalid vector returned: %v", r[1].GetData().(string))
	}
}

func TestDBGobRoundtrip(t *testing.T) {
	a, err := NewVector([]float64{10, 9}, "kris is a potato")
	if err != nil {
		t.Fatal(err)
	}
	b, err := NewVector([]float64{-5, 0}, "john is a fool")
	if err != nil {
		t.Fatal(err)
	}
	c, err := NewVector([]float64{5000, 10000}, "fritz is a german")
	if err != nil {
		t.Fatal(err)
	}

	db := new(DB)
	db.SetParallel(2)
	db.SetFilename("test.db")
	db.AddVector(a)
	db.AddVector(b)
	db.AddVector(c)

	encoded, err := db.GobEncode()
	if err != nil {
		t.Fatal(err)
	}

	db2 := new(DB)
	if err := db2.GobDecode(encoded); err != nil {
		t.Fatal(err)
	}

	if db2.Len() != 3 {
		t.Fatalf("expected 3 vectors, got %d", db2.Len())
	}

	// verify fields
	if db2.filename != "test.db" {
		t.Fatalf("expected filename 'test.db', got %q", db2.filename)
	}
	if !db2.fileBacked {
		t.Fatal("expected fileBacked to be true")
	}
	if db2.p != 2 {
		t.Fatalf("expected p=2, got %d", db2.p)
	}

	// verify vector data
	expected := []string{"kris is a potato", "john is a fool", "fritz is a german"}
	for i, v := range db2.vectors {
		if v.GetData().(string) != expected[i] {
			t.Fatalf("vector %d: expected %q, got %q", i, expected[i], v.GetData().(string))
		}
	}

	// verify linked list (prev/next)
	if db2.vectors[0].Prev() != nil {
		t.Fatal("first vector prev should be nil")
	}
	if db2.vectors[0].Next() != db2.vectors[1] {
		t.Fatal("first vector next should be second vector")
	}
	if db2.vectors[1].Prev() != db2.vectors[0] {
		t.Fatal("second vector prev should be first vector")
	}
	if db2.vectors[2].Next() != nil {
		t.Fatal("last vector next should be nil")
	}

	// verify the decoded DB still works for similarity queries
	x, err := NewVector([]float64{8, 9}, "query")
	if err != nil {
		t.Fatal(err)
	}
	r, err := db2.SimilarVectors(x, 1, 0.1)
	if err != nil {
		t.Fatal(err)
	}
	if r[0].GetData().(string) != "kris is a potato" {
		t.Fatalf("similarity query: expected 'kris is a potato', got %q", r[0].GetData().(string))
	}
}

func TestDBParallel(t *testing.T) {
	db := new(DB)
	db.SetParallel(31) // should leave a remainder

	for i := 0; i < 1000; i++ {
		x, err := NewVector([]float64{float64(i)}, fmt.Sprintf("data%v", i))
		if err != nil {
			t.Fatal(err)
		}
		db.AddVector(x)
	}
	if db.Len() != 1000 {
		t.Fatal("invalid length")
	}

	x, err := NewVector([]float64{0}, "Who is a potato?")
	if err != nil {
		t.Fatal(err)
	}

	r, err := db.SimilarVectors(x, 1000, -1)
	if err != nil {
		t.Fatal(err)
	}
	if len(r) != 1000 {
		t.Fatal("invalid response", len(r))
	}
}

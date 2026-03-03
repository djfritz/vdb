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

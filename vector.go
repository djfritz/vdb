package vdb

import (
	"encoding/json"
	"fmt"
	"math"
	"slices"
)

type Vector struct {
	V    []float64 `json:"v"`
	N    float64   `json:"n"`
	Data any       `json:"d"`
	prev *Vector
	next *Vector
}

type unmarshalVector struct {
	V    []float64       `json:"v"`
	N    float64         `json:"n"`
	Data json.RawMessage `json:"d"`
}

func (v *Vector) Prev() *Vector {
	return v.prev
}

func (v *Vector) Next() *Vector {
	return v.next
}

func (v *Vector) String() string {
	return fmt.Sprintf("%v", v.Data)
}

func NewVector(x []float64, data any) (*Vector, error) {
	if len(x) == 0 {
		return nil, ErrZeroVector
	}
	return &Vector{
		V:    x,
		N:    norm(x),
		Data: data,
	}, nil
}

func (v *Vector) Len() int {
	return len(v.V)
}

// Return the norm of the vector (the square root of the dot product of the vector
// with itself)
func norm(x []float64) float64 {
	var p float64
	for i := 0; i < len(x); i++ {
		p += x[i] * x[i]
	}
	return math.Sqrt(p)
}

func (v *Vector) GetData() any {
	return v.Data
}

func (a *Vector) Equal(b *Vector) bool {
	return slices.Equal(a.V, b.V)
}

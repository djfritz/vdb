package vdb

import (
	"errors"
)

var (
	ErrZeroVector         = errors.New("zero length vector")
	ErrVectorSizeMismatch = errors.New("vectors are not the same size")
)

// Finds the cosine similarity between two non-zero vectors.
func (a *Vector) CosineSimilarity(b *Vector) (float64, error) {
	if a.Len() == 0 || b.Len() == 0 {
		return 0, ErrZeroVector
	}
	if a.Len() != b.Len() {
		return 0, ErrVectorSizeMismatch
	}

	var p float64
	for i := 0; i < a.Len(); i++ {
		p += a.V[i] * b.V[i]
	}
	s := p / (a.N * b.N)
	return s, nil
}

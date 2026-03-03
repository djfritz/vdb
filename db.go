package vdb

import (
	"container/heap"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"runtime"
	"slices"
	"sync"
)

var ErrInvalidParallel = errors.New("number of threads must be >= 0")

type DB struct {
	mtx        sync.Mutex
	vectors    []*Vector
	filename   string
	fileBacked bool
	p          int
}

func (d *DB) SetParallel(p int) error {
	if p < 0 {
		return ErrInvalidParallel
	}
	d.p = p
	return nil
}

func (d *DB) SetFilename(f string) {
	d.mtx.Lock()
	defer d.mtx.Unlock()
	d.filename = f
	if f == "" {
		d.fileBacked = false
	} else {
		d.fileBacked = true
	}
}

type DecodeFunction func(json.RawMessage) (any, error)

func (d *DB) LoadFile(fn DecodeFunction) error {
	d.mtx.Lock()
	defer d.mtx.Unlock()
	if d.filename == "" {
		return nil
	}

	_, err := os.Stat(d.filename)
	if err != nil {
		return err
	}

	f, err := os.Open(d.filename)
	if err != nil {
		return err
	}
	defer f.Close()

	dec := json.NewDecoder(f)
	for {
		var x *unmarshalVector
		err := dec.Decode(&x)
		if err != nil {
			if err == io.EOF {
				return nil
			}
			return err
		}

		v := &Vector{
			V: x.V,
			N: x.N,
		}
		ddata, err := fn(x.Data)
		if err != nil {
			return err
		}
		v.Data = ddata

		if len(d.vectors) != 0 {
			v.prev = d.vectors[len(d.vectors)-1]
			d.vectors[len(d.vectors)-1].next = v
		}
		d.vectors = append(d.vectors, v)
	}
	return nil
}

func (d *DB) AddVector(x *Vector) error {
	d.mtx.Lock()
	defer d.mtx.Unlock()
	for _, v := range d.vectors {
		if v.Equal(x) {
			return nil
		}
	}
	if len(d.vectors) != 0 {
		x.prev = d.vectors[len(d.vectors)-1]
		d.vectors[len(d.vectors)-1].next = x
	}
	d.vectors = append(d.vectors, x)
	return d.commit(x)
}

func (d *DB) commit(x *Vector) error {
	if !d.fileBacked {
		return nil
	}

	f, err := os.OpenFile(d.filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	return enc.Encode(x)
}

type similarVector struct {
	s float64
	v *Vector
}

type similarityHeap []*similarVector

func (s similarityHeap) Len() int           { return len(s) }
func (s similarityHeap) Less(i, j int) bool { return s[i].s < s[j].s }
func (s similarityHeap) Swap(i, j int)      { s[i], s[j] = s[j], s[i] }

func (s *similarityHeap) Push(x any) {
	*s = append(*s, x.(*similarVector))
}

func (s *similarityHeap) Pop() any {
	old := *s
	n := len(old)
	x := old[n-1]
	*s = old[0 : n-1]
	return x
}

func (d *DB) Len() int {
	d.mtx.Lock()
	defer d.mtx.Unlock()
	return len(d.vectors)
}

func (d *DB) SimilarVectors(x *Vector, n int, t float64) ([]*Vector, error) {
	d.mtx.Lock()
	defer d.mtx.Unlock()
	h := new(similarityHeap)
	heap.Init(h)

	p := d.p
	if p == 0 {
		p = runtime.NumCPU()
	}

	// split vectors in to p partitions, or the length of vectors if it's less than p
	part := p
	if len(d.vectors) < p {
		part = len(d.vectors)
	}
	step := len(d.vectors) / part
	remainder := len(d.vectors) % part

	coalesce := make(chan *similarVector, len(d.vectors))
	var wg sync.WaitGroup
	var last int
	for i := range part {
		wg.Add(1)
		start := last
		end := start + step
		if i < remainder {
			end++
		}
		fmt.Println(start, end)
		if end > len(d.vectors) {
			end = len(d.vectors)
		}
		go similarPartition(coalesce, x, d.vectors[start:end], t, &wg)
		last = end
	}

	wg.Wait()
	close(coalesce)

	for s := range coalesce {
		heap.Push(h, s)
		if h.Len() > n {
			heap.Pop(h)
		}
	}

	var r []*Vector
	for h.Len() > 0 {
		r = append(r, heap.Pop(h).(*similarVector).v)
	}
	slices.Reverse(r)

	return r, nil
}

func similarPartition(coalesce chan *similarVector, x *Vector, vectors []*Vector, t float64, wg *sync.WaitGroup) {
	defer wg.Done()
	for _, v := range vectors {
		cs, err := x.CosineSimilarity(v)
		if err != nil {
			panic(err)
		}

		if cs < t {
			continue
		}

		s := &similarVector{
			s: cs,
			v: v,
		}
		coalesce <- s
	}
}

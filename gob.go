// Copyright 2025 David Fritz. All rights reserved.
// This software may be modified and distributed under the terms of the BSD
// 2-clause license. See the LICENSE file for details.

package vdb

import (
	"bytes"
	"encoding/gob"
)

func (d *DB) GobEncode() ([]byte, error) {
	d.mtx.Lock()
	defer d.mtx.Unlock()

	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)

	if err := enc.Encode(d.filename); err != nil {
		return nil, err
	}
	if err := enc.Encode(d.fileBacked); err != nil {
		return nil, err
	}
	if err := enc.Encode(d.p); err != nil {
		return nil, err
	}
	if err := enc.Encode(len(d.vectors)); err != nil {
		return nil, err
	}
	for _, v := range d.vectors {
		if err := enc.Encode(v.V); err != nil {
			return nil, err
		}
		if err := enc.Encode(v.N); err != nil {
			return nil, err
		}
		if err := enc.Encode(&v.Data); err != nil {
			return nil, err
		}
	}

	return buf.Bytes(), nil
}

func (d *DB) GobDecode(data []byte) error {
	dec := gob.NewDecoder(bytes.NewReader(data))

	if err := dec.Decode(&d.filename); err != nil {
		return err
	}
	if err := dec.Decode(&d.fileBacked); err != nil {
		return err
	}
	if err := dec.Decode(&d.p); err != nil {
		return err
	}

	var n int
	if err := dec.Decode(&n); err != nil {
		return err
	}

	d.vectors = make([]*Vector, 0, n)
	for i := range n {
		v := &Vector{}
		if err := dec.Decode(&v.V); err != nil {
			return err
		}
		if err := dec.Decode(&v.N); err != nil {
			return err
		}
		if err := dec.Decode(&v.Data); err != nil {
			return err
		}
		if i > 0 {
			v.prev = d.vectors[i-1]
			d.vectors[i-1].next = v
		}
		d.vectors = append(d.vectors, v)
	}

	return nil
}

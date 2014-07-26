package leveldb

// #cgo LDFLAGS: -lleveldb
// #include <stdlib.h>
// #include "leveldb/c.h"
// #include "leveldb_ext.h"
import "C"

import (
	"bytes"
	"unsafe"
)

const (
	IteratorForward  uint8 = 0
	IteratorBackward uint8 = 1
)

const (
	RangeClose uint8 = 0x00
	RangeLOpen uint8 = 0x01
	RangeROpen uint8 = 0x10
	RangeOpen  uint8 = 0x11
)

// min must less or equal than max
//
// range type:
//
// 	close: [min, max]
// 	open: (min, max)
// 	lopen: (min, max]
// 	ropen: [min, max)
//
type Range struct {
	Min []byte
	Max []byte

	Type uint8
}

type Limit struct {
	Offset int
	Count  int
}

type Iterator struct {
	it      *C.leveldb_iterator_t
	isValid C.uchar
}

// Returns a copy of key.
func (it *Iterator) Key() []byte {
	var klen C.size_t
	kdata := C.leveldb_iter_key(it.it, &klen)
	if kdata == nil {
		return nil
	}

	return C.GoBytes(unsafe.Pointer(kdata), C.int(klen))
}

// Returns a copy of value.
func (it *Iterator) Value() []byte {
	var vlen C.size_t
	vdata := C.leveldb_iter_value(it.it, &vlen)
	if vdata == nil {
		return nil
	}

	return C.GoBytes(unsafe.Pointer(vdata), C.int(vlen))
}

// Returns a reference of key.
// you must be careful that it will be changed after next iterate.
func (it *Iterator) RawKey() []byte {
	var klen C.size_t
	kdata := C.leveldb_iter_key(it.it, &klen)
	if kdata == nil {
		return nil
	}

	return slice(unsafe.Pointer(kdata), int(C.int(klen)))
}

// Returns a reference of value.
// you must be careful that it will be changed after next iterate.
func (it *Iterator) RawValue() []byte {
	var vlen C.size_t
	vdata := C.leveldb_iter_value(it.it, &vlen)
	if vdata == nil {
		return nil
	}

	return slice(unsafe.Pointer(vdata), int(C.int(vlen)))
}

// Copy key to b, if b len is small or nil, returns a new one.
func (it *Iterator) BufKey(b []byte) []byte {
	k := it.RawKey()
	if k == nil {
		return nil
	}
	if b == nil {
		b = []byte{}
	}

	b = b[0:0]
	return append(b, k...)
}

// Copy value to b, if b len is small or nil, returns a new one.
func (it *Iterator) BufValue(b []byte) []byte {
	v := it.RawValue()
	if v == nil {
		return nil
	}

	if b == nil {
		b = []byte{}
	}

	b = b[0:0]
	return append(b, v...)
}

func (it *Iterator) Close() {
	if it.it != nil {
		C.leveldb_iter_destroy(it.it)
		it.it = nil
	}
}

func (it *Iterator) Valid() bool {
	return ucharToBool(it.isValid)
}

func (it *Iterator) Next() {
	it.isValid = C.leveldb_iter_next_ext(it.it)
}

func (it *Iterator) Prev() {
	it.isValid = C.leveldb_iter_prev_ext(it.it)
}

func (it *Iterator) SeekToFirst() {
	it.isValid = C.leveldb_iter_seek_to_first_ext(it.it)
}

func (it *Iterator) SeekToLast() {
	it.isValid = C.leveldb_iter_seek_to_last_ext(it.it)
}

func (it *Iterator) Seek(key []byte) {
	it.isValid = C.leveldb_iter_seek_ext(it.it, (*C.char)(unsafe.Pointer(&key[0])), C.size_t(len(key)))
}

// Finds by key, if not found, nil returns.
func (it *Iterator) Find(key []byte) []byte {
	it.Seek(key)
	if it.Valid() {
		k := it.RawKey()
		if k == nil {
			return nil
		} else if bytes.Equal(k, key) {
			return it.Value()
		}
	}

	return nil
}

// Finds by key, if not found, nil returns, else a reference of value returns.
// you must be careful that it will be changed after next iterate.
func (it *Iterator) RawFind(key []byte) []byte {
	it.Seek(key)
	if it.Valid() {
		k := it.RawKey()
		if k == nil {
			return nil
		} else if bytes.Equal(k, key) {
			return it.RawValue()
		}
	}

	return nil
}

type RangeLimitIterator struct {
	it *Iterator

	r *Range
	l *Limit

	step int

	//0 for IteratorForward, 1 for IteratorBackward
	direction uint8
}

func (it *RangeLimitIterator) Key() []byte {
	return it.it.Key()
}

func (it *RangeLimitIterator) Value() []byte {
	return it.it.Value()
}

func (it *RangeLimitIterator) RawKey() []byte {
	return it.it.RawKey()
}

func (it *RangeLimitIterator) RawValue() []byte {
	return it.it.RawValue()
}

func (it *RangeLimitIterator) BufKey(b []byte) []byte {
	return it.it.BufKey(b)
}

func (it *RangeLimitIterator) BufValue(b []byte) []byte {
	return it.it.BufValue(b)
}

func (it *RangeLimitIterator) Valid() bool {
	if it.l.Offset < 0 {
		return false
	} else if !it.it.Valid() {
		return false
	} else if it.l.Count >= 0 && it.step >= it.l.Count {
		return false
	}

	if it.direction == IteratorForward {
		if it.r.Max != nil {
			r := bytes.Compare(it.it.RawKey(), it.r.Max)
			if it.r.Type&RangeROpen > 0 {
				return !(r >= 0)
			} else {
				return !(r > 0)
			}
		}
	} else {
		if it.r.Min != nil {
			r := bytes.Compare(it.it.RawKey(), it.r.Min)
			if it.r.Type&RangeLOpen > 0 {
				return !(r <= 0)
			} else {
				return !(r < 0)
			}
		}
	}

	return true
}

func (it *RangeLimitIterator) Next() {
	it.step++

	if it.direction == IteratorForward {
		it.it.Next()
	} else {
		it.it.Prev()
	}
}

func (it *RangeLimitIterator) Close() {
	it.it.Close()
}

func NewRangeLimitIterator(i *Iterator, r *Range, l *Limit) *RangeLimitIterator {
	return rangeLimitIterator(i, r, l, IteratorForward)
}

func NewRevRangeLimitIterator(i *Iterator, r *Range, l *Limit) *RangeLimitIterator {
	return rangeLimitIterator(i, r, l, IteratorBackward)
}

func NewRangeIterator(i *Iterator, r *Range) *RangeLimitIterator {
	return rangeLimitIterator(i, r, &Limit{0, -1}, IteratorForward)
}

func NewRevRangeIterator(i *Iterator, r *Range) *RangeLimitIterator {
	return rangeLimitIterator(i, r, &Limit{0, -1}, IteratorBackward)
}

func rangeLimitIterator(i *Iterator, r *Range, l *Limit, direction uint8) *RangeLimitIterator {
	it := new(RangeLimitIterator)

	it.it = i

	it.r = r
	it.l = l
	it.direction = direction

	it.step = 0

	if l.Offset < 0 {
		return it
	}

	if direction == IteratorForward {
		if r.Min == nil {
			it.it.SeekToFirst()
		} else {
			it.it.Seek(r.Min)

			if r.Type&RangeLOpen > 0 {
				if it.it.Valid() && bytes.Equal(it.it.RawKey(), r.Min) {
					it.it.Next()
				}
			}
		}
	} else {
		if r.Max == nil {
			it.it.SeekToLast()
		} else {
			it.it.Seek(r.Max)

			if !it.it.Valid() {
				it.it.SeekToLast()
			} else {
				if !bytes.Equal(it.it.RawKey(), r.Max) {
					it.it.Prev()
				}
			}

			if r.Type&RangeROpen > 0 {
				if it.it.Valid() && bytes.Equal(it.it.RawKey(), r.Max) {
					it.it.Prev()
				}
			}
		}
	}

	for i := 0; i < l.Offset; i++ {
		if it.it.Valid() {
			if it.direction == IteratorForward {
				it.it.Next()
			} else {
				it.it.Prev()
			}
		}
	}

	return it
}

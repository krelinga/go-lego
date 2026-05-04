package pod

import (
	"iter"

	"github.com/krelinga/go-libs/tuple"
)

// SliceDict is a Dict implementation that preserves the order of insertion.
// Iteration over the dict will yield entries in the order they were added.
// Key equality is determined by a user-provided equality function, allowing
// for custom comparison logic.
//
// Many of the operations on SliceDict require iterating through all entries
// to check for key equality, resulting in O(n) time complexity for operations like
// Get, Put, and Del. This makes SliceDict less efficient than a hash-based dict
// for large collections, but it provides the benefit of maintaining insertion
// order and allowing for custom key equality logic.
type SliceDict[K, V any] struct {
	slice    *Slice[tuple.T2[K, V]]
	equalKey func(a, b K) bool
}

// NewSliceDictFunc creates a new SliceDict with a custom key equality function and initial entries.
// If any of the initial entries have keys that are considered equal, only the last occurrence will be kept.
func NewSliceDictFunc[K, V any](equalKey func(a, b K) bool, entries ...tuple.T2[K, V]) *SliceDict[K, V] {
	d := &SliceDict[K, V]{
		slice:    &Slice[tuple.T2[K, V]]{},
		equalKey: equalKey,
	}
	d.Reserve(len(entries))
	for _, e := range entries {
		d.Put(e.A, e.B)
	}
	return d
}

// NewSliceDict creates a new SliceDict using Go's built-in equality for keys and the given initial entries.
// If any of the initial entries have equal keys, only the last occurrence will be kept.
func NewSliceDict[K comparable, V any](entries ...tuple.T2[K, V]) *SliceDict[K, V] {
	return NewSliceDictFunc(func(a, b K) bool { return a == b }, entries...)
}

// NewSliceDictOfFunc creates a new SliceDict from a Bag2, using a custom key equality function.
// The dict will contain all unique-key entries from the Bag2, in the order they appear.
// The new SliceDict will be independent of the original Bag2.
func NewSliceDictOfFunc[K, V any](equalKey func(a, b K) bool, b Bag2[K, V]) *SliceDict[K, V] {
	d := NewSliceDictFunc[K, V](equalKey)
	d.Reserve(b.Len())
	for k, v := range b.Elems() {
		d.Put(k, v)
	}
	return d
}

// NewSliceDictOf creates a new SliceDict from a Bag2, using Go's built-in equality for keys.
// The new SliceDict will be independent of the original Bag2.
func NewSliceDictOf[K comparable, V any](b Bag2[K, V]) *SliceDict[K, V] {
	return NewSliceDictOfFunc(func(a, b K) bool { return a == b }, b)
}

// findKey returns the index of the entry with the given key, or -1 if not found.
func (d *SliceDict[K, V]) findKey(key K) int {
	for i := range d.slice.Len() {
		if d.equalKey(d.slice.Get(i).A, key) {
			return i
		}
	}
	return -1
}

// Len returns the number of entries in the SliceDict.
func (d *SliceDict[K, V]) Len() int {
	return d.slice.Len()
}

// Get retrieves the value associated with the given key.
// It returns the value and true if the key was found, or the zero value and false otherwise.
// This operation has O(n) time complexity.
func (d *SliceDict[K, V]) Get(key K) (V, bool) {
	i := d.findKey(key)
	if i < 0 {
		var zero V
		return zero, false
	}
	return d.slice.Get(i).B, true
}

// KeyVals returns a sequence of key-value pairs in the order they were added.
func (d *SliceDict[K, V]) KeyVals() iter.Seq2[K, V] {
	return func(yield func(K, V) bool) {
		for e := range d.slice.Vals() {
			if !yield(e.A, e.B) {
				return
			}
		}
	}
}

// Keys returns a sequence of keys in the order they were added.
func (d *SliceDict[K, V]) Keys() iter.Seq[K] {
	return func(yield func(K) bool) {
		for e := range d.slice.Vals() {
			if !yield(e.A) {
				return
			}
		}
	}
}

// Vals returns a sequence of values in the order they were added.
func (d *SliceDict[K, V]) Vals() iter.Seq[V] {
	return func(yield func(V) bool) {
		for e := range d.slice.Vals() {
			if !yield(e.B) {
				return
			}
		}
	}
}

// Put adds or updates the entry for the given key.
// If the key already exists, the existing entry is removed from the dict and the new entry is added at the end, preserving insertion order.
// This operation has O(n) time complexity due to the need to search for the key.
func (d *SliceDict[K, V]) Put(key K, value V) {
	d.Del(key) // Remove existing entry if it exists to maintain insertion order
	d.slice.Push(tuple.New2(key, value))
}

// Clear removes all entries from the SliceDict.
func (d *SliceDict[K, V]) Clear() {
	d.slice.Clear()
}

// Del removes the entry with the given key from the SliceDict.
// If the key is not found, this operation has no effect.
// This operation has O(n) time complexity.
func (d *SliceDict[K, V]) Del(key K) {
	i := d.findKey(key)
	if i < 0 {
		return
	}
	for j := i + 1; j < d.slice.Len(); j++ {
		d.slice.Set(j-1, d.slice.Get(j))
	}
	d.slice.Pop()
}

// Reserve pre-allocates space for at least n entries in the SliceDict.
func (d *SliceDict[K, V]) Reserve(n int) {
	d.slice.Reserve(n)
}

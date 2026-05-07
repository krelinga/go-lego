// Package pod provides a collection of data structures and utilities for working with sets, dictionaries, and vectors.
//
// It includes interfaces and implementations for mutable and immutable collections, as well as helper functions for creating and manipulating these collections.
package pod

// Viewer is implemented by types that can produce a view of themselves as type T.
type Viewer[T any] interface {
	View() T
}

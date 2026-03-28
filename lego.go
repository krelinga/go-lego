/*
Package lego provides several abstractions on top of Go's built-in data structures, such as maps and slices. It defines interfaces and functions that allow you to work with these data structures in a more flexible and composable way.

# Key Concepts

## Maps, Sets, and Slices

The package provides wrappers around Go's built-in maps and slices, and implements a set type that wraps a map[K]struct{}.

## Pairs

The package uses Pair types to represent key-value pairs in maps and elements in slices (where keys are slice indices).

## Fixed Types

A Fixed Type provides a view of another container that does not allow adding or removing elements.
The elements in the container may still be modified, for example for containers of pointers.

Fixed types are expressed as interfaces: [FixedMap], [FixedSet], and [FixedSlice].

## View Types

A View Type provides an immutable view of another type.
The elements of the type (if any) must be immutable themselves.

It is common for a mutable type to have a corresponding view type, possibly with slightly modified types for members.
Such mutable types can implement the [Viewer] interface, which allows them to provide a view of themselves.

This package defines [ViewerMap], [ViewerSet], and [ViewerSlice] types that can store values that implement [Viewer] and provide views of them.

# Caveats

## Mutable Set Values and Map Keys

This package assumes that all map keys and set values are immutable.
This assumption breaks down, for example, if map keys or set values are pointers to mutable types.
It is advised to only use immutable types as map keys and set values when using this package, and to avoid modifying map keys and set values after they have been added to a map or set.
*/
package lego
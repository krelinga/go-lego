package match

import (
	"fmt"
	"reflect"
	"runtime"
)

type Meta struct {
	Name       string
	SourceFile string
	SourceLine int
}

func (m Meta) String() string {
	// TODO: include source file and line number.
	return m.Name
}

func MetaHere() Meta {
	pc, file, line, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()

	return Meta{
		Name:       funcName,
		SourceFile: file,
		SourceLine: line,
	}
}

type Result struct {
	// Information about the matcher itself.
	Meta Meta

	// Value that was matched against.
	Val reflect.Value

	// Did the [Matcher] accept Val?
	Accepted bool

	// Human-readable explanation of the Result.
	Why string

	// Any other values that are useful to understand the context of the Result.
	Context []Context

	// Results of any child matchers.  This is used to build a tree of match results, which can be used for debugging and error reporting.
	Children []Child
}

type Child struct {
	// Name of the child in the context of the parent matcher.
	//
	// This may be empty if there is no meaningful name for the child.
	Name string

	// Result from the child matcher.
	Result *Result
}

type Context struct {
	// Name of the context value.
	Name string

	// Value of the context.
	Val reflect.Value
}

type FatalError struct {
	// Information about the matcher that failed.
	Meta Meta
	
	// Value that was matched against.
	Val reflect.Value

	// The fatal error.
	Err error
}

func (e *FatalError) Error() string {
	return fmt.Sprintf("FATAL: %s: %s", e.Meta, e.Err.Error())
}

func (e *FatalError) Unwrap() error {
	return e.Err
}

type ChildError struct {
	// Name of the child in the context of the parent matcher.
	//
	// This may be empty if there is no meaningful name for the child.
	Name string

	// The error from the child matcher.
	Err error
}

func (e *ChildError) Error() string {
	namePart := ""
	if e.Name != "" {
		namePart = fmt.Sprintf(" %s", e.Name)
	}
	return fmt.Sprintf("child%s: %s", namePart, e.Err.Error())
}

func (e *ChildError) Unwrap() error {
	return e.Err
}
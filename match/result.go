package match

import (
	"reflect"
	"runtime"
)

type Meta struct {
	Name string
	SourceFile string
	SourceLine int
}

func (m Meta) WithName(name string) Meta {
	m.Name = name
	return m
}

func MetaHere() Meta {
	pc, file, line, _ := runtime.Caller(1)
	funcName := runtime.FuncForPC(pc).Name()
	
	return Meta{
		Name: funcName,
		SourceFile: file,
		SourceLine: line,
	}
}

type Result struct {
	// Information about the matcher itself.
	Meta Meta

	// Value that was matched against.
	Val reflect.Value

	// Did the matcher accept Val?  Or was there an error?
	Accepted bool
	Err error

	// A function that can generate a human-readable description of the match result.
	Explain func() string

	// Results of any child matchers.  This is used to build a tree of match results, which can be used for debugging and error reporting.
	Children []Child
}

type Child struct {
	// Name of the child in the context of the parent matcher.
	Name string

	// Result from the child matcher.
	Result *Result
}

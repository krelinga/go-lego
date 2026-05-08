package match

import (
	"reflect"
	"runtime"
)

type Meta struct {
	Name       string
	SourceFile string
	SourceLine int
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

	// Did the [Matcher] accept Val?  Or was there an error?
	Accepted bool
	Err      error

	// A function that can generate a human-readable description of the match result.
	//
	// If this is nil then a default explanation will be generated.
	Explain func() string

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

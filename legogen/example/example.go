package main

import "github.com/krelinga/go-lego"

type Example struct {
	StrField string
	IntField int
}

func (e *Example) View() ExampleView {
	return ExampleView{e}
}

type ExampleView struct {
	e *Example
}

func (v ExampleView) StrField() string {
	return v.e.StrField
}

func (v ExampleView) IntField() int {
	return v.e.IntField
}

func (v ExampleView) Compare(other ExampleView) int {
	return lego.CompareUsing(v, other,
		lego.NewCmpFuncGo(ExampleView.StrField),
		lego.NewCmpFuncGo(ExampleView.IntField),
	)
}

//go:generate go run github.com/krelinga/go-lego/legogen -type=*Example -equal=viewer
//go:generate go run github.com/krelinga/go-lego/legogen -type=ExampleView -equal=comparer

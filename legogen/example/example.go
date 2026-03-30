package main

import "github.com/krelinga/go-lego"

type Example struct {
	lego.GenView
	lego.GenEqualComparer
	StrField string `legogen:"view=getter"`
	IntField int    `legogen:"view=getter"`
}
//go:generate go run github.com/krelinga/go-lego/legogen -type=Example

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

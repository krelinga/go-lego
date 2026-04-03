package main

import (
	"fmt"

	v2 "github.com/krelinga/go-lego/v2"
)

type SKU string

type SKUCounts struct {
	v2.Map[SKU, int]
}

type SKUCountsView = v2.FixedDict[SKU, int]

func (s *SKUCounts) View() SKUCountsView {
	return s
}

func NewSKUCounts(kvs ...v2.KV[SKU, int]) *SKUCounts {
	m := SKUCounts{}
	v2.AddAll(&m.Map, kvs...)
	return &m
}

type Location string

type Inventory struct {
	v2.Map[Location, *SKUCounts]
}

func (i *Inventory) View() InventoryView {
	return InventoryView{
		DictViewEmbed: v2.NewMapViewEmbed(i),
	}
}

func NewInventory(kvs ...v2.KV[Location, *SKUCounts]) *Inventory {
	m := Inventory{}
	v2.AddAll(&m.Map, kvs...)
	return &m
}

type InventoryView struct {
	v2.DictViewEmbed[Location, *SKUCounts, SKUCountsView]
}

func main() {
	inv := NewInventory(
		v2.NewKV(Location("North America"), NewSKUCounts(
			v2.NewKV(SKU("Widget"), 100),
			v2.NewKV(SKU("Gizmo"), 50),
		)),
	)
	fmt.Println(inv.Length())
}

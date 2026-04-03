package main

import (
	"fmt"

	v2 "github.com/krelinga/go-lego/v2"
)

type SKU string

type SKUCounts struct {
	v2.Map[SKU, int]
}

func (s *SKUCounts) View() SKUCountsView {
	return SKUCountsView{
		FixedDict: s,
		DictEqualEmbedComparable: v2.NewDictEqualEmbedComparable[SKUCountsView](s),
	}
}

func (s *SKUCounts) Total() int {
	return s.View().Total()
}

func (s SKUCounts) Equal(other SKUCounts) bool {
	return s.View().Equal(other.View())
}

func NewSKUCounts(kvs ...v2.KV[SKU, int]) *SKUCounts {
	m := SKUCounts{}
	v2.AddAll(&m.Map, kvs...)
	return &m
}

type SKUCountsView struct {
	v2.FixedDict[SKU, int]
	v2.DictEqualEmbedComparable[SKUCountsView, SKU, int]
}

func (s SKUCountsView) Total() int {
	total := 0
	for count := range s.Values() {
		total += count
	}
	return total
}

type Location string

type Inventory struct {
	v2.Map[Location, *SKUCounts]
}

func (i *Inventory) View() InventoryView {
	dve := v2.NewMapViewEmbed(i)
	return InventoryView{
		DictViewEmbed: dve,
		DictEqualEmbed: v2.NewDictEqualEmbed[InventoryView](dve),
	}
}

func (i *Inventory) Total() int {
	return i.View().Total()
}

func (i Inventory) Equal(other *Inventory) bool {
	return i.View().Equal(other.View())
}

func NewInventory(kvs ...v2.KV[Location, *SKUCounts]) *Inventory {
	m := Inventory{}
	v2.AddAll(&m.Map, kvs...)
	return &m
}

type InventoryView struct {
	v2.DictViewEmbed[Location, *SKUCounts, SKUCountsView]
	v2.DictEqualEmbed[InventoryView, Location, SKUCountsView]
}

func (v InventoryView) Total() int {
	total := 0
	for skuCounts := range v.Values() {
		total += skuCounts.Total()
	}
	return total
}

func main() {
	inv := NewInventory(
		v2.NewKV(Location("North America"), NewSKUCounts(
			v2.NewKV(SKU("Widget"), 100),
			v2.NewKV(SKU("Gizmo"), 50),
		)),
	)
	fmt.Println(inv.View().Total())
	fmt.Println(inv.Total())
	fmt.Println(inv)

	inv2 := NewInventory()
	fmt.Println(inv.Equal(inv2))
	fmt.Println(inv.View().Equal(inv2.View()))
	fmt.Println(inv.Equal(inv))
	fmt.Println(inv.View().Equal(inv.View()))
}

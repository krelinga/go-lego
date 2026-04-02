package main

import v2 "github.com/krelinga/go-lego/v2"

type SKU string

type SKUCounts struct {
	v2.Map[SKU, int]
}

type SKUCountsView = v2.FixedMap[SKU, int]

func (s SKUCounts) View() SKUCountsView {
	return s
}

type Location string

type Inventory struct {
	v2.Map[Location, SKUCounts]
}

func (i Inventory) View() InventoryView {
	return InventoryView{
		MapViewEmbed: v2.NewMapVieEmbed(i),
	}
}

type InventoryView struct {
	v2.MapViewEmbed[Location, SKUCounts, SKUCountsView]
}
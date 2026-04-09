package main

import (
	"fmt"

	"github.com/krelinga/go-lego/pod"
)

///////////////////////////////////////////////////////////////////////////////

type SKU string

///////////////////////////////////////////////////////////////////////////////

//pod:container
//pod:view name=SKUCountsView values=direct
//pod:equal values=comparable
type SKUCounts struct {
	pod.Map[SKU, int]
}

func (s *SKUCounts) View() SKUCountsView {
	return SKUCountsView{
		FixedDict: s,
	}
}

//pod:view expose=direct
func (s *SKUCounts) Total() int {
	return s.View().Total()
}

func (s SKUCounts) Equal(other SKUCounts) bool {
	return s.View().Equal(other.View())
}

func NewSKUCounts(kvs ...pod.KV[SKU, int]) *SKUCounts {
	m := SKUCounts{}
	pod.AddAll(&m.Map, pod.RangeFrom(kvs...))
	return &m
}

type SKUCountsView struct {
	pod.FixedDict[SKU, int]
}

func (s SKUCountsView) Equal(other SKUCountsView) bool {
	return pod.DictEqualValuesComparable(s, other)
}

func (s SKUCountsView) Total() int {
	total := 0
	for count := range s.Values() {
		total += count
	}
	return total
}

///////////////////////////////////////////////////////////////////////////////

type Location string

///////////////////////////////////////////////////////////////////////////////

//pod:container
//pod:view name=InventoryView values=SKUCountsView
//pod:equal values=Equal
type Inventory struct {
	pod.Map[Location, *SKUCounts]
}

func (i *Inventory) View() InventoryView {
	dve := pod.NewDictViewEmbed(i)
	return InventoryView{
		DictViewEmbed: dve,
	}
}

//pod:view expose=direct
func (i *Inventory) Total() int {
	return i.View().Total()
}

func (i Inventory) Equal(other *Inventory) bool {
	return i.View().Equal(other.View())
}

func NewInventory(kvs ...pod.KV[Location, *SKUCounts]) *Inventory {
	m := Inventory{}
	pod.AddAll(&m.Map, pod.RangeFrom(kvs...))
	return &m
}

type InventoryView struct {
	pod.DictViewEmbed[Location, *SKUCounts, SKUCountsView]
}

func (v InventoryView) Equal(other InventoryView) bool {
	return pod.DictEqualValues(v, other)
}

func (v InventoryView) Total() int {
	total := 0
	for skuCounts := range v.Values() {
		total += skuCounts.Total()
	}
	return total
}

///////////////////////////////////////////////////////////////////////////////

//pod:view name=StatusCountsView
//pod:compare order=desc(Ready),Backordered,WaitingToShip
//pod:equal use=Compare
type StatusCounts struct {
	Ready         int
	Backordered   int
	WaitingToShip int
}

func (s *StatusCounts) View() StatusCountsView {
	return StatusCountsView{
		sc: s,
	}
}

func (s *StatusCounts) Compare(other *StatusCounts) int {
	return s.View().Compare(other.View())
}

func (s StatusCounts) Equal(other *StatusCounts) bool {
	return s.View().Equal(other.View())
}

type StatusCountsView struct {
	sc *StatusCounts
}

func (s StatusCountsView) Ready() int {
	return s.sc.Ready
}

func (s StatusCountsView) Backordered() int {
	return s.sc.Backordered
}

func (s StatusCountsView) WaitingToShip() int {
	return s.sc.WaitingToShip
}

func (s StatusCountsView) Compare(other StatusCountsView) int {
	return pod.CompareUsing(s, other,
		pod.NewComparatorReversed(pod.NewComparatorOrdered(StatusCountsView.Ready)),
		pod.NewComparatorOrdered(StatusCountsView.Backordered),
		pod.NewComparatorOrdered(StatusCountsView.WaitingToShip),
	)
}

func (s StatusCountsView) Equal(other StatusCountsView) bool {
	return s.Compare(other) == 0
}

///////////////////////////////////////////////////////////////////////////////

//pod:container
//pod:view name=StatusCountsListView values=StatusCountsView
//pod:equal values=comparable
type StatusCountsList struct {
	pod.Slice[*StatusCounts]
}

func (l *StatusCountsList) View() StatusCountsListView {
	return StatusCountsListView{
		FixedList: l,
	}
}

func (l *StatusCountsList) Equal(other *StatusCountsList) bool {
	return l.View().Equal(other.View())
}

func NewStatusCountsList(items ...*StatusCounts) *StatusCountsList {
	sl := StatusCountsList{}
	pod.AddAll(&sl.Slice, pod.RangeFrom(items...))
	return &sl
}

type StatusCountsListView struct {
	pod.FixedList[int, *StatusCounts]
}

func (l StatusCountsListView) Equal(other StatusCountsListView) bool {
	return pod.ListEqualValues(l, other)
}

///////////////////////////////////////////////////////////////////////////////

func main() {
	inv := NewInventory(
		pod.NewKV(Location("North America"), NewSKUCounts(
			pod.NewKV(SKU("Widget"), 100),
			pod.NewKV(SKU("Gizmo"), 50),
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

	sc1 := &StatusCounts{Ready: 10, Backordered: 5, WaitingToShip: 2}
	sc2 := &StatusCounts{Ready: 8, Backordered: 10, WaitingToShip: 1}
	fmt.Println("pod.LessThan(sc1, sc2):", pod.LessThan(sc1, sc2))
	fmt.Println("pod.GreaterThan(sc1, sc2):", pod.GreaterThan(sc1, sc2))
	fmt.Println("sc1.Equal(sc2):", sc1.Equal(sc2))
	fmt.Println("sc1.Equal(sc1):", sc1.Equal(sc1))

	scl1 := NewStatusCountsList(sc1, sc2)
	scl2 := NewStatusCountsList(sc2, sc1)
	fmt.Println("scl1.Equal(scl2):", scl1.Equal(scl2))
	fmt.Println("scl1.View().Equal(scl2.View()):", scl1.View().Equal(scl2.View()))
	fmt.Println("scl1.Equal(scl1):", scl1.Equal(scl1))
	fmt.Println("scl1.View().Equal(scl1.View()):", scl1.View().Equal(scl1.View()))
}

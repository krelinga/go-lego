package conx

import (
	"iter"
	"maps"
)

type DictView[K, V any] interface {
	Len() int
	Get(key K) (V, bool)
	All() iter.Seq2[K, V]
}

type Dict[K comparable, V any] struct {
	data map[K]V
}

func NewDict[K comparable, V any](data ...KV[K, V]) (d *Dict[K, V]) {
	d = &Dict[K, V]{}
	if len(data) == 0 {
		return
	}
	d.Reserve(len(data))
	for _, kv := range data {
		d.Set(kv.K, kv.V)
	}
	return
}

func CloneDict[K comparable, V any](dict DictView[K, V]) *Dict[K, V] {
	return CloneDictFunc(dict, func(k K) K { return k }, func(v V) V { return v })
}

func CloneDictFunc[K any, KK comparable, V, VV any](dict DictView[K, V], keyFunc func(K) KK, valueFunc func(V) VV) *Dict[KK, VV] {
	d := &Dict[KK, VV]{}
	d.Reserve(dict.Len())
	for k, v := range dict.All() {
		d.Set(keyFunc(k), valueFunc(v))
	}
	return d
}

func (d Dict[K, V]) Len() int {
	return len(d.data)
}

func (d Dict[K, V]) Get(key K) (V, bool) {
	value, ok := d.data[key]
	return value, ok
}

func (d Dict[K, V]) All() iter.Seq2[K, V] {
	return maps.All(d.data)
}

func (d *Dict[K, V]) Set(key K, value V) {
	if d.data == nil {
		d.data = make(map[K]V)
	}
	d.data[key] = value
}

func (d *Dict[K, V]) Clear() {
	d.data = nil
}

func (d *Dict[K, V]) Reserve(n int) {
	if d.data == nil {
		d.data = make(map[K]V, n)
	}
}

func (d *Dict[K, V]) Delete(key K) {
	delete(d.data, key)
}

func WrapDictValues[K, V, W any](dict DictView[K, V], wrap func(V) W) DictView[K, W] {
	return wrappedDictValues[K, V, W]{
		dict: dict,
		wrap: wrap,
	}
}

type wrappedDictValues[K, V, W any] struct {
	dict DictView[K, V]
	wrap func(V) W
}

func (w wrappedDictValues[K, V, W]) Len() int {
	return w.dict.Len()
}

func (w wrappedDictValues[K, V, W]) Get(key K) (W, bool) {
	value, ok := w.dict.Get(key)
	if !ok {
		var zero W
		return zero, false
	}
	return w.wrap(value), true
}

func (w wrappedDictValues[K, V, W]) All() iter.Seq2[K, W] {
	return func(yield func(K, W) bool) {
		for k, v := range w.dict.All() {
			if !yield(k, w.wrap(v)) {
				return
			}
		}
	}
}

func WrapDictKeys[K, V, W any](dict DictView[K, V], wrap func(K) W, unwrap func(W) K) DictView[W, V] {
	return wrappedDictKeys[K, V, W]{
		dict:   dict,
		wrap:   wrap,
		unwrap: unwrap,
	}
}

type wrappedDictKeys[K, V, W any] struct {
	dict   DictView[K, V]
	wrap   func(K) W
	unwrap func(W) K
}

func (w wrappedDictKeys[K, V, W]) Len() int {
	return w.dict.Len()
}

func (w wrappedDictKeys[K, V, W]) Get(key W) (V, bool) {
	return w.dict.Get(w.unwrap(key))
}

func (w wrappedDictKeys[K, V, W]) All() iter.Seq2[W, V] {
	return func(yield func(W, V) bool) {
		for k, v := range w.dict.All() {
			if !yield(w.wrap(k), v) {
				return
			}
		}
	}
}

func NewDictViewFromMap[M ~map[K]V, K comparable, V any](m M) DictView[K, V] {
	return mapDictView[M, K, V]{m: m}
}

type mapDictView[M ~map[K]V, K comparable, V any] struct {
	m M
}

func (v mapDictView[M, K, V]) Len() int {
	return len(v.m)
}

func (v mapDictView[M, K, V]) Get(key K) (V, bool) {
	value, ok := v.m[key]
	return value, ok
}

func (v mapDictView[M, K, V]) All() iter.Seq2[K, V] {
	return maps.All(v.m)
}

func DictAll[K, V any](dict DictView[K, V]) Range2[K, V] {
	return dictAllRange[K, V]{dict: dict}
}

type dictAllRange[K, V any] struct {
	dict DictView[K, V]
}

func (r dictAllRange[K, V]) Len() int {
	return r.dict.Len()
}

func (r dictAllRange[K, V]) Range() iter.Seq2[K, V] {
	return r.dict.All()
}

func DictKeys[K, V any](dict DictView[K, V]) Range[K] {
	return dictKeysRange[K, V]{dict: dict}
}

type dictKeysRange[K, V any] struct {
	dict DictView[K, V]
}

func (r dictKeysRange[K, V]) Len() int {
	return r.dict.Len()
}

func (r dictKeysRange[K, V]) Range() iter.Seq[K] {
	return func(yield func(K) bool) {
		for k := range r.dict.All() {
			if !yield(k) {
				return
			}
		}
	}
}

func DictValues[K, V any](dict DictView[K, V]) Range[V] {
	return dictValuesRange[K, V]{dict: dict}
}

type dictValuesRange[K, V any] struct {
	dict DictView[K, V]
}

func (r dictValuesRange[K, V]) Len() int {
	return r.dict.Len()
}

func (r dictValuesRange[K, V]) Range() iter.Seq[V] {
	return func(yield func(V) bool) {
		for _, v := range r.dict.All() {
			if !yield(v) {
				return
			}
		}
	}
}

func DictEqual[K, V comparable](a, b DictView[K, V]) bool {
	return DictEqualFunc(a, b, func(vA, vB V) bool {
		return vA == vB
	})
}

func DictEqualFunc[K, V any](a, b DictView[K, V], eq func(V, V) bool) bool {
	if a.Len() != b.Len() {
		return false
	}
	for k, vA := range a.All() {
		vB, ok := b.Get(k)
		if !ok || !eq(vA, vB) {
			return false
		}
	}
	return true
}

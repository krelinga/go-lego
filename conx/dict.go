package conx

type Dict[K comparable, V any] struct {
	data map[K]V
}

func NewDict[K comparable, V any](data ...KV[K, V]) *Dict[K, V] {
	if len(data) == 0 {
		return &Dict[K, V]{}
	}
	m := make(map[K]V, len(data))
	for _, kv := range data {
		m[kv.K] = kv.V
	}
	return &Dict[K, V]{data: m}

}

func (d *Dict[K, V]) Len() int {
	return len(d.data)
}

func (d *Dict[K, V]) Find(key K) (V, bool) {
	v, ok := d.data[key]
	return v, ok
}

func (d *Dict[K, V]) Set(key K, value V) {
	d.data[key] = value
}

func (d *Dict[K, V]) Delete(key K) {
	delete(d.data, key)
}

func (d *Dict[K, V]) Clear() {
	d.data = nil
}

func (d *Dict[K, V]) Reserve(capacity int) {
	if d.data == nil {
		d.data = make(map[K]V, capacity)
	}
}
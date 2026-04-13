package conx

type KV[K, V any] struct {
	K K
	V V
}

func NewKV[K, V any](k K, v V) KV[K, V] {
	return KV[K, V]{K: k, V: v}
}

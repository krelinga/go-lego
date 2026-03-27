package lego

type Lener interface {
	Len() int
}

func Len[L Lener](l L) int {
	return l.Len()
}
package v2

type CanLength interface {
	Length() int
}

func Length(s CanLength) int {
	return s.Length()
}

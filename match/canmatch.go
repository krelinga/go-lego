package match

type CanMatch interface {
	Match(any, Reporter) bool
}

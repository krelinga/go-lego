package order

// Func is an ordering function for values of type T.
// It returns a negative value if x < y, zero if x == y, and a positive value if x > y.
// This convention matches the signature used by cmp.Compare.
type Func[T any] = func(x T, y T) int

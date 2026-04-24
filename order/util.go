package order

// Less reports whether the ordering result indicates x < y.
func Less(order int) bool {
	return order < 0
}

// Greater reports whether the ordering result indicates x > y.
func Greater(order int) bool {
	return order > 0
}

// Equal reports whether the ordering result indicates x == y.
func Equal(order int) bool {
	return order == 0
}

// LessEqual reports whether the ordering result indicates x <= y.
func LessEqual(order int) bool {
	return order <= 0
}

// GreaterEqual reports whether the ordering result indicates x >= y.
func GreaterEqual(order int) bool {
	return order >= 0
}

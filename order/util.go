package order

func Less(order int) bool {
	return order < 0
}

func Greater(order int) bool {
	return order > 0
}

func Equal(order int) bool {
	return order == 0
}

func LessEqual(order int) bool {
	return order <= 0
}

func GreaterEqual(order int) bool {
	return order >= 0
}
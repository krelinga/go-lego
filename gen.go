package lego

// GenView is an empty struct that can be embedded in a struct to indicate that it should have a generated view type.
type GenView struct{}

// GenEqualComparer is an empty struct that can be embedded in a struct to indicate that it should have a generated Equal method based on its Compare method.
type GenEqualComparer struct{}
package helpers

import "constraints"

func EqualValues[V constraints.Ordered](a, b *V) bool {
	if a != nil && b != nil {
		return *a == *b
	}
	if a == nil && b == nil {
		return true
	}
	return false
}

func CopyValue[V constraints.Ordered](from *V) *V {
	var copy *V
	if from == nil {
		return nil
	}
	*copy = *from
	return copy
}

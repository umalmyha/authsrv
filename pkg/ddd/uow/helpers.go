package uow

import (
	"reflect"
)

func isPointer(data any) bool {
	return reflect.ValueOf(data).Kind() == reflect.Ptr
}

func isNilPtr(data any) bool {
	v := reflect.ValueOf(data)
	if v.Kind() == reflect.Ptr && v.IsNil() {
		return true
	}
	return false
}

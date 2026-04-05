package tsv

import "reflect"

func isComplexType(typ reflect.Type) bool {
	switch typ.Kind() {
	case reflect.Struct, reflect.Slice, reflect.Array, reflect.Map, reflect.Ptr:
		return true
	default:
		return false
	}
}

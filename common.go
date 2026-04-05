package tsv

import "reflect"

func isComplexType(typ reflect.Type) bool {
	switch typ.Kind() {
	case reflect.Struct, reflect.Map:
		return true
	case reflect.Pointer:
		return isComplexType(typ.Elem())
	default:
		return false
	}
}

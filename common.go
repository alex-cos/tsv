package tsv

import (
	"bytes"
	"reflect"
)

func isMapType(typ reflect.Type) bool {
	if typ.Kind() == reflect.Map {
		return true
	}
	if typ.Kind() == reflect.Ptr {
		return isMapType(typ.Elem())
	}
	return false
}

func escapeString(s, delimiter string) string {
	var buf bytes.Buffer
	start := 0
	for i := range len(s) {
		var repl string
		switch s[i] {
		case '\\':
			repl = `\\`
		case '\t':
			repl = `\t`
		case '\n':
			repl = `\n`
		case '\r':
			repl = `\r`
		default:
			if delimiter != "" && s[i] == delimiter[0] {
				repl = ` `
			} else {
				continue
			}
		}
		buf.WriteString(s[start:i])
		buf.WriteString(repl)
		start = i + 1
	}
	buf.WriteString(s[start:])
	return buf.String()
}

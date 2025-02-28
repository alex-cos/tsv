package tsv

import (
	"bytes"
	"fmt"
	"reflect"
	"strconv"
)

const (
	tab = 0x09
	eol = 0x0A
)

type encoderFunc func(v reflect.Value) ([]byte, error)

// Encoder defines struture for TSV Encoder.
type Encoder struct {
	nice       bool
	timeFormat string
}

// NewTSVEncoder builds and return a new TSVEncoder.
func NewTSVEncoder(nice bool) *Encoder {
	return &Encoder{
		nice:       nice,
		timeFormat: "2006/01/02 15:04:05",
	}
}

// Encode encodes given interface to TSV format.
func (e *Encoder) Encode(v any) ([]byte, error) {
	var buf bytes.Buffer

	val := reflect.ValueOf(v)
	encoder := e.typeEncoder(val.Type())

	b, err := encoder(val)
	if err != nil {
		return buf.Bytes(), err
	}
	buf.Write(b)

	return buf.Bytes(), nil
}

func (e *Encoder) typeEncoder(typ reflect.Type) encoderFunc {
	switch typ.Kind() {
	case reflect.Bool:
		return boolEncoder
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return intEncoder
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64, reflect.Uintptr:
		return uintEncoder
	case reflect.Float32:
		return float32Encoder
	case reflect.Float64:
		return float64Encoder
	case reflect.String:
		return stringEncoder
	case reflect.Interface:
		return e.interfaceEncoder
	case reflect.Struct:
		name := fmt.Sprintf("%s.%s", typ.PkgPath(), typ.Name())
		if name == "time.Time" {
			if e.nice {
				return e.timeNiceEncoder
			}
			return timeEncoder
		}
		return e.structEncoder
	case reflect.Map:
		return unsupportedTypeEncoder
	case reflect.Slice:
		return e.sliceEncoder
	case reflect.Array:
		return e.arrayEncoder
	case reflect.Ptr:
		return e.interfaceEncoder
	default:
		return unsupportedTypeEncoder
	}
}

func (e *Encoder) arrayEncoder(val reflect.Value) ([]byte, error) {
	var buf bytes.Buffer

	t := val.Type().Elem()
	encoder := e.typeEncoder(t)
	for i := 0; i < val.Len(); i++ {
		if i > 0 {
			k := t.Kind()
			if k == reflect.Array ||
				k == reflect.Slice ||
				k == reflect.Struct {
				buf.WriteByte(eol)
			} else {
				buf.WriteByte(tab)
			}
		}
		b, err := encoder(val.Index(i))
		if err != nil {
			return buf.Bytes(), err
		}
		buf.Write(b)
	}
	return buf.Bytes(), nil
}

func (e *Encoder) interfaceEncoder(val reflect.Value) ([]byte, error) {
	var buf bytes.Buffer
	if val.IsNil() {
		buf.WriteString("null")
	} else {
		val := val.Elem()
		encoder := e.typeEncoder(val.Type())
		b, err := encoder(val)
		if err != nil {
			return buf.Bytes(), err
		}
		buf.Write(b)
	}
	return buf.Bytes(), nil
}

func (e *Encoder) sliceEncoder(val reflect.Value) ([]byte, error) {
	var buf bytes.Buffer
	if val.IsNil() {
		buf.WriteString("null")
	} else {
		b, err := e.arrayEncoder(val)
		if err != nil {
			return buf.Bytes(), err
		}
		buf.Write(b)
	}
	return buf.Bytes(), nil
}

func (e *Encoder) structEncoder(val reflect.Value) ([]byte, error) {
	var buf bytes.Buffer

	for i := 0; i < val.NumField(); i++ {
		if i > 0 {
			buf.WriteByte(tab)
		}
		f := val.Field(i)
		encoder := e.typeEncoder(f.Type())
		b, err := encoder(f)
		if err != nil {
			return buf.Bytes(), err
		}
		buf.Write(b)
	}
	return buf.Bytes(), nil
}

func boolEncoder(val reflect.Value) ([]byte, error) {
	var buf bytes.Buffer

	if val.Bool() {
		buf.WriteString("true")
	} else {
		buf.WriteString("false")
	}
	return buf.Bytes(), nil
}

func intEncoder(val reflect.Value) ([]byte, error) {
	b := []byte{}
	return strconv.AppendInt(b, val.Int(), 10), nil
}

func uintEncoder(val reflect.Value) ([]byte, error) {
	b := []byte{}
	return strconv.AppendUint(b, val.Uint(), 10), nil
}

func float32Encoder(val reflect.Value) ([]byte, error) {
	b := []byte{}
	return strconv.AppendFloat(b, val.Float(), 'f', -1, 32), nil
}

func float64Encoder(val reflect.Value) ([]byte, error) {
	b := []byte{}
	return strconv.AppendFloat(b, val.Float(), 'f', -1, 64), nil
}

func stringEncoder(val reflect.Value) ([]byte, error) {
	var buf bytes.Buffer
	buf.WriteString(val.String())
	return buf.Bytes(), nil
}

func timeEncoder(val reflect.Value) ([]byte, error) {
	_, ok := val.Type().MethodByName("Unix")
	if !ok {
		return nil, fmt.Errorf("wrong time value: %s", val.Type().String())
	}
	method := val.MethodByName("Unix")
	res := method.Call([]reflect.Value{})
	if len(res) != 1 {
		return nil, fmt.Errorf("wrong time value: %s", val.Type().String())
	}
	return intEncoder(res[0])
}

func (e *Encoder) timeNiceEncoder(val reflect.Value) ([]byte, error) {
	_, ok := val.Type().MethodByName("Format")
	if !ok {
		return nil, fmt.Errorf("wrong time value: %s", val.Type().String())
	}
	method := val.MethodByName("Format")
	in := []reflect.Value{reflect.ValueOf(e.timeFormat)}
	res := method.Call(in)
	if len(res) != 1 {
		return nil, fmt.Errorf("wrong time value: %s", val.Type().String())
	}
	return stringEncoder(res[0])
}

func unsupportedTypeEncoder(val reflect.Value) ([]byte, error) {
	return nil, fmt.Errorf("unsupported type: %s", val.Type().String())
}

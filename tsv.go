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

type encoderFunc func(buf *bytes.Buffer, v reflect.Value) error

// Encoder defines structure for TSV Encoder.
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

	if err := encoder(&buf, val); err != nil {
		return buf.Bytes(), err
	}

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
		return e.interfaceEncoderFn()
	case reflect.Struct:
		name := fmt.Sprintf("%s.%s", typ.PkgPath(), typ.Name())
		if name == "time.Time" {
			if e.nice {
				return e.timeNiceEncoderFn()
			}
			return timeEncoder
		}
		return e.structEncoderFn(typ)
	case reflect.Map:
		return unsupportedTypeEncoder
	case reflect.Slice:
		return e.sliceEncoderFn(typ)
	case reflect.Array:
		return e.arrayEncoderFn(typ)
	case reflect.Ptr:
		return e.interfaceEncoderFn()
	default:
		return unsupportedTypeEncoder
	}
}

func (e *Encoder) arrayEncoderFn(typ reflect.Type) encoderFunc {
	elemType := typ.Elem()
	encoder := e.typeEncoder(elemType)

	return func(buf *bytes.Buffer, val reflect.Value) error {
		for i := range val.Len() {
			if i > 0 {
				k := elemType.Kind()
				if k == reflect.Array ||
					k == reflect.Slice ||
					k == reflect.Struct {
					buf.WriteByte(eol)
				} else {
					buf.WriteByte(tab)
				}
			}
			if err := encoder(buf, val.Index(i)); err != nil {
				return err
			}
		}
		return nil
	}
}

func (e *Encoder) interfaceEncoderFn() encoderFunc {
	return func(buf *bytes.Buffer, val reflect.Value) error {
		if val.IsNil() {
			buf.WriteString("null")
		} else {
			val := val.Elem()
			encoder := e.typeEncoder(val.Type())
			return encoder(buf, val)
		}
		return nil
	}
}

func (e *Encoder) sliceEncoderFn(typ reflect.Type) encoderFunc {
	enc := e.arrayEncoderFn(typ)

	return func(buf *bytes.Buffer, val reflect.Value) error {
		if val.IsNil() {
			buf.WriteString("null")
			return nil
		}
		return enc(buf, val)
	}
}

func (e *Encoder) structEncoderFn(typ reflect.Type) encoderFunc {
	n := typ.NumField()
	type fieldInfo struct {
		encoder encoderFunc
		index   int
	}
	fields := make([]fieldInfo, n)
	for i := range n {
		fields[i] = fieldInfo{
			encoder: e.typeEncoder(typ.Field(i).Type),
			index:   i,
		}
	}

	return func(buf *bytes.Buffer, val reflect.Value) error {
		for i, fi := range fields {
			if i > 0 {
				buf.WriteByte(tab)
			}
			if err := fi.encoder(buf, val.Field(fi.index)); err != nil {
				return err
			}
		}
		return nil
	}
}

func boolEncoder(buf *bytes.Buffer, val reflect.Value) error {
	if val.Bool() {
		buf.WriteString("true")
	} else {
		buf.WriteString("false")
	}
	return nil
}

func intEncoder(buf *bytes.Buffer, val reflect.Value) error {
	buf.Write(strconv.AppendInt(nil, val.Int(), 10))
	return nil
}

func uintEncoder(buf *bytes.Buffer, val reflect.Value) error {
	buf.Write(strconv.AppendUint(nil, val.Uint(), 10))
	return nil
}

func float32Encoder(buf *bytes.Buffer, val reflect.Value) error {
	buf.Write(strconv.AppendFloat(nil, val.Float(), 'f', -1, 32))
	return nil
}

func float64Encoder(buf *bytes.Buffer, val reflect.Value) error {
	buf.Write(strconv.AppendFloat(nil, val.Float(), 'f', -1, 64))
	return nil
}

func stringEncoder(buf *bytes.Buffer, val reflect.Value) error {
	buf.WriteString(val.String())
	return nil
}

func timeEncoder(buf *bytes.Buffer, val reflect.Value) error {
	_, ok := val.Type().MethodByName("Unix")
	if !ok {
		return fmt.Errorf("wrong time value: %s", val.Type().String())
	}
	method := val.MethodByName("Unix")
	res := method.Call([]reflect.Value{})
	if len(res) != 1 {
		return fmt.Errorf("wrong time value: %s", val.Type().String())
	}
	return intEncoder(buf, res[0])
}

func (e *Encoder) timeNiceEncoderFn() encoderFunc {
	return func(buf *bytes.Buffer, val reflect.Value) error {
		_, ok := val.Type().MethodByName("Format")
		if !ok {
			return fmt.Errorf("wrong time value: %s", val.Type().String())
		}
		method := val.MethodByName("Format")
		in := []reflect.Value{reflect.ValueOf(e.timeFormat)}
		res := method.Call(in)
		if len(res) != 1 {
			return fmt.Errorf("wrong time value: %s", val.Type().String())
		}
		return stringEncoder(buf, res[0])
	}
}

func unsupportedTypeEncoder(buf *bytes.Buffer, val reflect.Value) error {
	return fmt.Errorf("unsupported type: %s", val.Type().String())
}

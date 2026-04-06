package tsv

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"reflect"
	"strconv"
	"time"
)

const (
	tab = 0x09
	cr  = 0x0D
	lf  = 0x0A
)

type encoderFunc func(buf *bytes.Buffer, v reflect.Value) error

// Encoder defines structure for TSV Encoder.
type Encoder struct {
	timeFormat string
	utc        bool
	crlf       bool
	delimiter  rune
}

// NewTSVEncoder builds and returns a new TSVEncoder.
// By default, time.Time values are encoded as Unix epoch timestamps.
// Use WithTimeFormat to specify a custom date/time format.
func NewTSVEncoder(opts ...Option) *Encoder {
	e := &Encoder{
		timeFormat: "",
		utc:        false,
		crlf:       false,
		delimiter:  0x09,
	}
	for _, opt := range opts {
		opt(e)
	}

	return e
}

func (e *Encoder) delim() string {
	if e.delimiter != 0 {
		return string(e.delimiter)
	}
	return string(rune(tab))
}

func (e *Encoder) endln() []byte {
	if e.crlf {
		return []byte{cr, lf}
	} else {
		return []byte{lf}
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

// EncodeTo encodes given interface to TSV format and writes it to writer.
func (e *Encoder) EncodeTo(w io.Writer, v any) error {
	b, err := e.Encode(v)
	if err != nil {
		return err
	}
	_, err = w.Write(b)
	return err
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
		return e.stringEncoderFn()
	case reflect.Interface:
		return e.interfaceEncoderFn()
	case reflect.Struct:
		if typ == reflect.TypeOf(time.Time{}) {
			return e.timeEncoder
		}
		return e.structEncoderFn(typ)
	case reflect.Map:
		return e.mapEncoderFn(typ)
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

	if isComplexType(elemType) {
		return e.complexEncoderFn(elemType)
	}

	encoder := e.typeEncoder(elemType)
	kind := elemType.Kind()
	return func(buf *bytes.Buffer, val reflect.Value) error {
		for i := range val.Len() {
			if i > 0 {
				if kind == reflect.Array ||
					kind == reflect.Slice {
					buf.Write(e.endln())
				} else {
					buf.WriteString(e.delim())
				}
			}
			if err := encoder(buf, val.Index(i)); err != nil {
				return err
			}
		}
		return nil
	}
}

func (e *Encoder) complexEncoderFn(elemType reflect.Type) encoderFunc {
	return func(buf *bytes.Buffer, val reflect.Value) error {
		for i := range val.Len() {
			if i > 0 {
				buf.WriteString(e.delim())
			}
			elem := val.Index(i)
			if elemType.Kind() == reflect.Ptr && elem.IsNil() {
				continue
			}
			b, err := json.Marshal(elem.Interface())
			if err != nil {
				return err
			}
			buf.Write(b)
		}
		return nil
	}
}

func (e *Encoder) interfaceEncoderFn() encoderFunc {
	return func(buf *bytes.Buffer, val reflect.Value) error {
		if val.IsNil() {
			return nil
		}
		elem := val.Elem()
		encoder := e.typeEncoder(elem.Type())
		return encoder(buf, elem)
	}
}

func (e *Encoder) sliceEncoderFn(typ reflect.Type) encoderFunc {
	enc := e.arrayEncoderFn(typ)

	return func(buf *bytes.Buffer, val reflect.Value) error {
		if val.IsNil() {
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
				buf.WriteString(e.delim())
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

func (e *Encoder) stringEncoderFn() encoderFunc {
	return func(buf *bytes.Buffer, val reflect.Value) error {
		s := val.String()
		del := e.delim()
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
				if del != "" && s[i] == del[0] {
					repl = `\` + del
				} else {
					continue
				}
			}
			buf.WriteString(s[start:i])
			buf.WriteString(repl)
			start = i + 1
		}
		buf.WriteString(s[start:])
		return nil
	}
}

func (e *Encoder) timeEncoder(buf *bytes.Buffer, val reflect.Value) error {
	t, ok := val.Interface().(time.Time)
	if !ok {
		return nil
	}
	if e.timeFormat != "" {
		buf.WriteString(t.Format(e.timeFormat))
	} else {
		if e.utc {
			t = t.UTC()
		}
		buf.Write(strconv.AppendInt(nil, t.Unix(), 10))
	}
	return nil
}

func (e *Encoder) mapEncoderFn(typ reflect.Type) encoderFunc {
	keyEncoder := e.typeEncoder(typ.Key())
	valueEncoder := e.typeEncoder(typ.Elem())

	return func(buf *bytes.Buffer, val reflect.Value) error {
		if val.IsNil() {
			return nil
		}

		first := true
		iter := val.MapRange()
		for iter.Next() {
			if !first {
				buf.Write(e.endln())
			}
			if err := keyEncoder(buf, iter.Key()); err != nil {
				return err
			}
			buf.WriteString(e.delim())
			if err := valueEncoder(buf, iter.Value()); err != nil {
				return err
			}
			first = false
		}
		return nil
	}
}

func unsupportedTypeEncoder(buf *bytes.Buffer, val reflect.Value) error {
	return fmt.Errorf("unsupported type: %s", val.Type().String())
}

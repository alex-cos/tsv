package tsv_test

import (
	"bytes"
	"encoding/json"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/alex-cos/tsv"
	"github.com/stretchr/testify/assert"
)

func TestBool(t *testing.T) {
	t.Parallel()

	boolean, err := tsv.NewTSVEncoder().Encode(true)
	assert.NoError(t, err)
	assert.Equal(t, "true", string(boolean))

	boolean, err = tsv.NewTSVEncoder().Encode(false)
	assert.NoError(t, err)
	assert.Equal(t, "false", string(boolean))
}

func TestString(t *testing.T) {
	t.Parallel()

	b, err := tsv.NewTSVEncoder().Encode("Test")
	assert.NoError(t, err)
	assert.Equal(t, "Test", string(b))
}

func TestInt(t *testing.T) {
	t.Parallel()

	input := int(32)
	b, err := tsv.NewTSVEncoder().Encode(input)
	assert.NoError(t, err)
	assert.Equal(t, "32", string(b))
}

func TestUInt(t *testing.T) {
	t.Parallel()

	input := uint(32)
	b, err := tsv.NewTSVEncoder().Encode(input)
	assert.NoError(t, err)
	assert.Equal(t, "32", string(b))
}

func TestFloat32(t *testing.T) {
	t.Parallel()

	input := float32(2789.98)
	b, err := tsv.NewTSVEncoder().Encode(input)
	assert.NoError(t, err)
	assert.Equal(t, "2789.98", string(b))
}

func TestFloat64(t *testing.T) {
	t.Parallel()

	input := float64(2789.9801)
	b, err := tsv.NewTSVEncoder().Encode(input)
	assert.NoError(t, err)
	assert.Equal(t, "2789.9801", string(b))
}

func TestArray(t *testing.T) {
	var input [4]string

	t.Parallel()

	input[0] = "First"
	input[1] = "Second"
	input[2] = "Third"
	input[3] = "Fourth"
	b, err := tsv.NewTSVEncoder().Encode(input)
	assert.NoError(t, err)
	assert.Equal(t, "First\tSecond\tThird\tFourth", string(b))
}

func TestSlice(t *testing.T) {
	t.Parallel()

	input := []float64{1278.21, 907.9, 765.12, -12.87}
	b, err := tsv.NewTSVEncoder().Encode(input)
	assert.NoError(t, err)
	assert.Equal(t, "1278.21\t907.9\t765.12\t-12.87", string(b))
}

func TestArrayOfSlice(t *testing.T) {
	t.Parallel()

	var input [4][]int

	input[0] = []int{10, 11, 13, 14, 15}
	input[1] = []int{20, 21, 23, 24, 25}
	input[2] = []int{30, 31, 33, 34, 35}
	input[3] = []int{40, 41, 43, 44, 45}
	b, err := tsv.NewTSVEncoder().Encode(input)
	assert.NoError(t, err)
	assert.Equal(t, "10\t11\t13\t14\t15\n20\t21\t23\t24\t25\n30\t31\t33\t34\t35\n40\t41\t43\t44\t45", string(b))
}

func TestPointer(t *testing.T) {
	t.Parallel()

	input := "Test"
	b, err := tsv.NewTSVEncoder().Encode(&input)
	assert.NoError(t, err)
	assert.Equal(t, "Test", string(b))
}

func TestInterface(t *testing.T) {
	var input [1]interface{}

	t.Parallel()

	input[0] = "Test"
	b, err := tsv.NewTSVEncoder().Encode(input)
	assert.NoError(t, err)
	assert.Equal(t, "Test", string(b))
}

func TestStruct(t *testing.T) {
	t.Parallel()

	input := struct {
		v1 string
		v2 int
	}{
		v1: "Test",
		v2: 10,
	}
	res, err := tsv.NewTSVEncoder().Encode(input)
	assert.NoError(t, err)
	assert.Equal(t, "Test\t10", string(res))
}

func TestTimeEncoder(t *testing.T) {
	t.Parallel()

	input := time.Date(2020, 3, 23, 16, 24, 10, 0, time.UTC)
	res, err := tsv.NewTSVEncoder().Encode(input)
	assert.NoError(t, err)
	assert.Equal(t, "1584980650", string(res))

	res, err = tsv.NewTSVEncoder(tsv.WithTimeFormat("2006/01/02 15:04:05")).Encode(input)
	assert.NoError(t, err)
	assert.Equal(t, "2020/03/23 16:24:10", string(res))
}

func TestArrayOfStruct(t *testing.T) {
	t.Parallel()

	type row struct {
		V1 string `json:"v1"`
		V2 int    `json:"v2"`
	}
	input := []row{
		{V1: "Test1", V2: 10},
		{V1: "Test2", V2: 20},
		{V1: "Test3", V2: 30},
	}
	res, err := tsv.NewTSVEncoder().Encode(input)
	assert.NoError(t, err)
	assert.Equal(t, "Test1\t10\nTest2\t20\nTest3\t30", string(res))
}

func TestSliceOfStruct(t *testing.T) {
	t.Parallel()

	type row struct {
		V1 string `json:"v1"`
		V2 int    `json:"v2"`
	}
	input := []row{
		{V1: "A", V2: 1},
		{V1: "B", V2: 2},
	}
	res, err := tsv.NewTSVEncoder().Encode(input)
	assert.NoError(t, err)
	assert.Equal(t, "A\t1\nB\t2", string(res))
}

func TestSliceOfMap(t *testing.T) {
	t.Parallel()

	input := []map[string]int{
		{"a": 1, "b": 2},
		{"c": 3},
	}
	res, err := tsv.NewTSVEncoder().Encode(input)
	assert.NoError(t, err)

	lines := strings.Split(string(res), "\t")
	assert.Len(t, lines, 2)

	var m1, m2 map[string]int
	assert.NoError(t, json.Unmarshal([]byte(lines[0]), &m1))
	assert.Equal(t, 1, m1["a"])
	assert.Equal(t, 2, m1["b"])

	assert.NoError(t, json.Unmarshal([]byte(lines[1]), &m2))
	assert.Equal(t, 3, m2["c"])
}

func TestArrayOfMap(t *testing.T) {
	t.Parallel()

	var input [2]map[string]string
	input[0] = map[string]string{"x": "foo"}
	input[1] = map[string]string{"y": "bar"}
	res, err := tsv.NewTSVEncoder().Encode(input)
	assert.NoError(t, err)

	lines := strings.Split(string(res), "\t")
	assert.Len(t, lines, 2)

	var m1, m2 map[string]string
	assert.NoError(t, json.Unmarshal([]byte(lines[0]), &m1))
	assert.Equal(t, "foo", m1["x"])

	assert.NoError(t, json.Unmarshal([]byte(lines[1]), &m2))
	assert.Equal(t, "bar", m2["y"])
}

func TestNilPointerInSlice(t *testing.T) {
	t.Parallel()

	s1 := "hello"
	s3 := "world"
	input := []*string{&s1, nil, &s3}
	res, err := tsv.NewTSVEncoder().Encode(input)
	assert.NoError(t, err)
	assert.Equal(t, "hello\t\tworld", string(res))
}

func TestMapCRLF(t *testing.T) {
	t.Parallel()

	input := map[string]int{
		"a": 1,
		"b": 2,
	}
	res, err := tsv.NewTSVEncoder(tsv.WithCRLF()).Encode(input)
	assert.NoError(t, err)

	lines := strings.Split(string(res), "\r\n")
	assert.Len(t, lines, 2)

	expected := map[string]int{
		"a": 1,
		"b": 2,
	}
	for _, line := range lines {
		parts := strings.Split(line, "\t")
		assert.Len(t, parts, 2)
		key := parts[0]
		val, err := strconv.Atoi(parts[1])
		assert.NoError(t, err)
		assert.Equal(t, expected[key], val)
	}
}

func TestMap(t *testing.T) {
	t.Parallel()

	input := map[string]int{
		"v1": 10,
		"v2": 11,
		"v3": 12,
	}
	res, err := tsv.NewTSVEncoder().Encode(input)
	assert.NoError(t, err)

	lines := strings.Split(string(res), "\n")
	assert.Len(t, lines, 3)

	expected := map[string]int{
		"v1": 10,
		"v2": 11,
		"v3": 12,
	}
	for _, line := range lines {
		parts := strings.Split(line, "\t")
		assert.Len(t, parts, 2)
		key := parts[0]
		val, err := strconv.Atoi(parts[1])
		assert.NoError(t, err)
		assert.Equal(t, expected[key], val)
	}
}

func TestDelimiter(t *testing.T) {
	t.Parallel()

	type row struct {
		Name string
		Age  int
	}
	input := row{Name: "Alice", Age: 30}

	res, err := tsv.NewTSVEncoder(tsv.WithDelimiter(',')).Encode(input)
	assert.NoError(t, err)
	assert.Equal(t, "Alice,30", string(res))
}

func TestDelimiterInSlice(t *testing.T) {
	t.Parallel()

	input := []string{"foo", "bar", "baz"}

	res, err := tsv.NewTSVEncoder(tsv.WithDelimiter(';')).Encode(input)
	assert.NoError(t, err)
	assert.Equal(t, "foo;bar;baz", string(res))
}

func TestDelimiterEscaping(t *testing.T) {
	t.Parallel()

	input := []string{"hello;world", "foo"}

	res, err := tsv.NewTSVEncoder(tsv.WithDelimiter(';')).Encode(input)
	assert.NoError(t, err)
	assert.Equal(t, `hello world;foo`, string(res))
}

func TestEncodeTo(t *testing.T) {
	t.Parallel()

	input := []string{"foo", "bar", "baz"}

	var buf bytes.Buffer
	err := tsv.NewTSVEncoder().EncodeTo(&buf, input)
	assert.NoError(t, err)
	assert.Equal(t, "foo\tbar\tbaz", buf.String())
}

func TestEncodeToStruct(t *testing.T) {
	t.Parallel()

	type row struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}
	input := row{Name: "Alice", Age: 30}

	var buf bytes.Buffer
	err := tsv.NewTSVEncoder().EncodeTo(&buf, input)
	assert.NoError(t, err)
	assert.Equal(t, "Alice\t30", buf.String())
}

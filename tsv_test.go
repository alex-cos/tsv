package tsv_test

import (
	"testing"
	"time"

	"github.com/alex-cos/tsv"
	"github.com/stretchr/testify/assert"
)

func TestBool(t *testing.T) {
	t.Parallel()

	boolean, err := tsv.NewTSVEncoder(false).Encode(true)
	assert.NoError(t, err)
	assert.Equal(t, "true", string(boolean))

	boolean, err = tsv.NewTSVEncoder(false).Encode(false)
	assert.NoError(t, err)
	assert.Equal(t, "false", string(boolean))
}

func TestString(t *testing.T) {
	t.Parallel()

	b, err := tsv.NewTSVEncoder(false).Encode("Test")
	assert.NoError(t, err)
	assert.Equal(t, "Test", string(b))
}

func TestInt(t *testing.T) {
	t.Parallel()

	input := int(32)
	b, err := tsv.NewTSVEncoder(false).Encode(input)
	assert.NoError(t, err)
	assert.Equal(t, "32", string(b))
}

func TestUInt(t *testing.T) {
	t.Parallel()

	input := uint(32)
	b, err := tsv.NewTSVEncoder(false).Encode(input)
	assert.NoError(t, err)
	assert.Equal(t, "32", string(b))
}

func TestFloat32(t *testing.T) {
	t.Parallel()

	input := float32(2789.98)
	b, err := tsv.NewTSVEncoder(false).Encode(input)
	assert.NoError(t, err)
	assert.Equal(t, "2789.98", string(b))
}

func TestFloat64(t *testing.T) {
	t.Parallel()

	input := float64(2789.9801)
	b, err := tsv.NewTSVEncoder(false).Encode(input)
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
	b, err := tsv.NewTSVEncoder(false).Encode(input)
	assert.NoError(t, err)
	assert.Equal(t, "First\tSecond\tThird\tFourth", string(b))
}

func TestSlice(t *testing.T) {
	t.Parallel()

	input := []float64{1278.21, 907.9, 765.12, -12.87}
	b, err := tsv.NewTSVEncoder(false).Encode(input)
	assert.NoError(t, err)
	assert.Equal(t, "1278.21\t907.9\t765.12\t-12.87", string(b))
}

func TestArrayOfSlice(t *testing.T) {
	var input [4][]int

	t.Parallel()

	input[0] = []int{10, 11, 13, 14, 15}
	input[1] = []int{20, 21, 23, 24, 25}
	input[2] = []int{30, 31, 33, 34, 35}
	input[3] = []int{40, 41, 43, 44, 45}
	b, err := tsv.NewTSVEncoder(false).Encode(input)
	assert.NoError(t, err)
	assert.Equal(t, "10\t11\t13\t14\t15\n20\t21\t23\t24\t25\n30\t31\t33\t34\t35\n40\t41\t43\t44\t45", string(b))
}

func TestPointer(t *testing.T) {
	t.Parallel()

	input := "Test"
	b, err := tsv.NewTSVEncoder(false).Encode(&input)
	assert.NoError(t, err)
	assert.Equal(t, "Test", string(b))
}

func TestInterface(t *testing.T) {
	var input [1]interface{}

	t.Parallel()

	input[0] = "Test"
	b, err := tsv.NewTSVEncoder(false).Encode(input)
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
	res, err := tsv.NewTSVEncoder(false).Encode(input)
	assert.NoError(t, err)
	assert.Equal(t, "Test\t10", string(res))
}

func TestTimeEncoder(t *testing.T) {
	t.Parallel()

	input := time.Date(2020, 3, 23, 16, 24, 10, 0, time.UTC)
	res, err := tsv.NewTSVEncoder(false).Encode(input)
	assert.NoError(t, err)
	assert.Equal(t, "1584980650", string(res))

	res, err = tsv.NewTSVEncoder(true).Encode(input)
	assert.NoError(t, err)
	assert.Equal(t, "2020/03/23 16:24:10", string(res))
}

func TestArrayOfStruct(t *testing.T) {
	t.Parallel()

	input := []struct {
		v1 string
		v2 int
	}{
		{v1: "Test1", v2: 10},
		{v1: "Test2", v2: 20},
		{v1: "Test3", v2: 30},
	}
	res, err := tsv.NewTSVEncoder(false).Encode(input)
	assert.NoError(t, err)
	assert.Equal(t, "Test1\t10\nTest2\t20\nTest3\t30", string(res))
}

func TestMap(t *testing.T) {
	t.Parallel()

	input := map[string]int{
		"v1": 10,
		"v2": 11,
	}
	_, err := tsv.NewTSVEncoder(false).Encode(input)
	assert.Error(t, err)
}

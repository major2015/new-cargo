package log

import (
	"fmt"
	"reflect"
	"testing"

	"github.com/go-playground/assert/v2"
)

var prefix = "flatten_test"

func TestFlattenMap(t *testing.T) {
	result := make(FlatMap)

	data := make(map[string]interface{})
	data["foo"] = "bar"
	data["foo_2"] = "bar_2"

	flattenMap(result, prefix, reflect.ValueOf(data))

	assert.Equal(t, "bar", result[prefix+"."+"foo"])
	assert.Equal(t, "bar_2", result[prefix+"."+"foo_2"])

	batter := make(map[string]interface{})
	batter["A"] = 'A'
	batter["B"] = "B"
	data["map"] = batter

	flattenMap(result, prefix, reflect.ValueOf(data))

	assert.Equal(t, "B", result[prefix+"."+"map.B"])
	assert.Equal(t, "65", result[prefix+"."+"map.A"])
}

func TestFlattenSliceArray(t *testing.T) {
	result := make(FlatMap)

	slice := make([]int, 3)
	slice[0] = 0
	slice[1] = 1
	slice[2] = 2

	flattenSliceArray(result, prefix, reflect.ValueOf(slice))

	assert.Equal(t, "0", result[prefix+".0"])
	assert.Equal(t, "1", result[prefix+".1"])
	assert.Equal(t, "2", result[prefix+".2"])
}

type Batter struct {
	Name     string
	Price    int32
	Discount float32
}

func TestFlattenStruct(t *testing.T) {
	result := make(FlatMap)

	batter := Batter{
		Name:     "batter",
		Price:    10000,
		Discount: 0.85,
	}

	flattenStruct(result, prefix, reflect.ValueOf(batter))

	assert.Equal(t, "batter", result[prefix+".Name"])
	assert.Equal(t, "10000", result[prefix+".Price"])
	assert.Equal(t, "0.850000", result[prefix+".Discount"])
}

func TestFlatten(t *testing.T) {

	data := make(map[string]interface{})
	data["A"] = 'A'
	data["a"] = 'a'
	data["true"] = true
	data["false"] = false
	data["unknown"] = uint16(1)

	slice := make([]int, 3)
	slice[0] = 0
	slice[1] = 1
	slice[2] = 2
	data["slice"] = slice

	array := [3]string{}
	array[0] = "0"
	array[1] = "1"
	array[2] = "2"
	data["array"] = array

	batter := make(map[string]interface{})
	batter["foo"] = "bar"
	batter["foo_obj"] = Batter{
		Name:     "batter",
		Price:    10000,
		Discount: 0.85,
	}
	data["batter"] = batter

	result, err := Flatten(data)
	if err != nil {
		fmt.Errorf("testing error: %w", err)
		return
	}

	assert.Equal(t, "65", result["A"])
	assert.Equal(t, "97", result["a"])
	assert.Equal(t, "true", result["true"])
	assert.Equal(t, "false", result["false"])
	assert.Equal(t, "Unknown: %!s(uint16=1)", result["unknown"])

	assert.Equal(t, "0", result["slice.0"])
	assert.Equal(t, "1", result["slice.1"])
	assert.Equal(t, "2", result["slice.2"])

	assert.Equal(t, "0", result["slice.0"])
	assert.Equal(t, "1", result["slice.1"])
	assert.Equal(t, "2", result["slice.2"])

	assert.Equal(t, "bar", result["batter.foo"])
	assert.Equal(t, "batter", result["batter.foo_obj.Name"])
	assert.Equal(t, "10000", result["batter.foo_obj.Price"])
	assert.Equal(t, "0.850000", result["batter.foo_obj.Discount"])
}

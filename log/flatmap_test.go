package log

import (
	"reflect"
	"testing"

	"github.com/go-playground/assert/v2"
)

var prefix = "flatten_test"

func TestFlattenMap(t *testing.T) {
	result := make(FlatMap)

	data := make(map[string]string)
	data["foo"] = "bar"
	data["foo_2"] = "bar_2"

	flattenMap(result, prefix, reflect.ValueOf(data))

	assert.Equal(t, "bar", result[prefix+"."+"foo"])
	assert.Equal(t, "bar_2", result[prefix+"."+"foo_2"])
}

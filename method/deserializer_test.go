package method

import (
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestUnmarshaler(t *testing.T) {
	type s struct {
		Value string `json:"value"`
	}
	t.Run("unmarshal scalar", func(t *testing.T) {
		t.Parallel()
		var result interface{}
		err := Unmarshaler(&result, reflect.TypeOf(string(""))).UnmarshalJSON([]byte(`"value"`))
		require.NoError(t, err)
		assert.IsType(t, string(""), result)
		assert.Equal(t, "value", result)
	})
	t.Run("unmarshal struct", func(t *testing.T) {
		t.Parallel()
		var result interface{}
		err := Unmarshaler(&result, reflect.TypeOf(s{})).UnmarshalJSON([]byte(`{"value":"string"}`))
		require.NoError(t, err)
		assert.IsType(t, s{}, result)
		assert.Equal(t, s{"string"}, result)
	})
	t.Run("unmarshal struct ptr", func(t *testing.T) {
		t.Parallel()
		var result interface{}
		err := Unmarshaler(&result, reflect.TypeOf(&s{})).UnmarshalJSON([]byte(`{"value":"string"}`))
		require.NoError(t, err)
		assert.IsType(t, &s{}, result)
		assert.Equal(t, &s{"string"}, result)
	})
}

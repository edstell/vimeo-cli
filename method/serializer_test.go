package method

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestMarshaler(t *testing.T) {
	type s struct {
		Value string `json:"value"`
	}
	t.Run("marshal scalar", func(t *testing.T) {
		t.Parallel()
		result, err := Marshaler(reflect.ValueOf("string")).MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, []byte(`"string"`), result)
	})
	t.Run("marshal struct", func(t *testing.T) {
		t.Parallel()
		result, err := Marshaler(reflect.ValueOf(s{"string"})).MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, []byte(`{"value":"string"}`), result)
	})
	t.Run("marshal struct ptr", func(t *testing.T) {
		t.Parallel()
		result, err := Marshaler(reflect.ValueOf(&s{"string"})).MarshalJSON()
		require.NoError(t, err)
		assert.Equal(t, []byte(`{"value":"string"}`), result)
	})
}

func TestJSONSerializer(t *testing.T) {
	t.Parallel()
	var buf bytes.Buffer
	err := JSONSerializer(Marshaler).Serialize(&buf, []reflect.Value{
		reflect.ValueOf("a"),
		reflect.ValueOf("b"),
	})
	require.NoError(t, err)
	assert.Equal(t, `["a","b"]`+"\n", buf.String())
}

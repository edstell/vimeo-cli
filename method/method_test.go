package method

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type target struct{}

func (t *target) NoArgsNoReturn() {}

func (t *target) NoArgs() string {
	return "result"
}

func (t *target) FixedArgs(a string, b string) (string, string) {
	return a, b
}

func (t *target) VariadicArgs(a ...string) []string {
	return a
}

func TestCallerCall(t *testing.T) {
	t.Run("method with no arguments and no results", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		caller := NewCaller(reflect.ValueOf(&target{}))
		err := caller.Call("NoArgsNoReturn", nil, &buf)
		require.NoError(t, err)
		assert.Equal(t, `[]`+"\n", buf.String())
	})
	t.Run("method with no arguments", func(t *testing.T) {
		t.Parallel()
		var buf bytes.Buffer
		caller := NewCaller(reflect.ValueOf(&target{}))
		err := caller.Call("NoArgs", nil, &buf)
		require.NoError(t, err)
		assert.Equal(t, `["result"]`+"\n", buf.String())
	})
	t.Run("method with fixed arguments", func(t *testing.T) {
		t.Parallel()
		in := bytes.NewReader([]byte(`["a","b"]`))
		var out bytes.Buffer
		caller := NewCaller(reflect.ValueOf(&target{}))
		err := caller.Call("FixedArgs", in, &out)
		require.NoError(t, err)
		assert.Equal(t, `["a","b"]`+"\n", out.String())
	})
	t.Run("method with variadic arguments", func(t *testing.T) {
		t.Parallel()
		in := bytes.NewReader([]byte(`["a","b","c","d"]`))
		var out bytes.Buffer
		caller := NewCaller(reflect.ValueOf(&target{}))
		err := caller.Call("VariadicArgs", in, &out)
		require.NoError(t, err)
		assert.Equal(t, `[["a","b","c","d"]]`+"\n", out.String())
	})
}

package method

import (
	"encoding/json"
	"io"
	"reflect"
)

// Serializer implementations serialize values, writing the output to the
// io.Writer.
type Serializer interface {
	Serialize(io.Writer, []reflect.Value) error
}

type SerializerFunc func(io.Writer, []reflect.Value) error

// The SerializerFunc type is an adapter to allow the use of ordinary functions
// as serializers. If f is a function with the appropriate signature,
// SerializerFunc(f) is a Serializer that calls f.
func (f SerializerFunc) Serialize(w io.Writer, v []reflect.Value) error {
	return f(w, v)
}

// The MarshalerFunc type is an adapter to allow the use of ordinary
// functions as json.Marshaler[s]. If f is a function with the appropriate
// signature, MarshalerFunc(f) is a Marshaler that calls f.
type MarshalerFunc func() ([]byte, error)

func (f MarshalerFunc) MarshalJSON() ([]byte, error) {
	return f()
}

// Marshaler returns a json.Marshaler which will marshal the reflect.Value.
func Marshaler(v reflect.Value) json.Marshaler {
	return MarshalerFunc(func() ([]byte, error) {
		return json.Marshal(v.Interface())
	})
}

// JSONSerializer serializes the values to a json array.
func JSONSerializer(marshaler func(reflect.Value) json.Marshaler) Serializer {
	return SerializerFunc(func(w io.Writer, vs []reflect.Value) error {
		ms := make([]json.Marshaler, 0, len(vs))
		for _, v := range vs {
			ms = append(ms, marshaler(v))
		}
		return json.NewEncoder(w).Encode(ms)
	})
}

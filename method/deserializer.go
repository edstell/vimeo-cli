package method

import (
	"encoding/json"
	"errors"
	"io"
	"reflect"
)

// Deserializer implementations deserialize data from the io.Reader to a slice
// of reflect.Value[s].
type Deserializer interface {
	Deserialize(io.Reader, []reflect.Type) ([]reflect.Value, error)
}

// The DeserializerFunc type is an adapter to allow the use of ordinary
// functions as deserializes. If f is a function with the appropriate
// signature, DeserializerFunc(f) is a Deserializer that calls f.
type DeserializerFunc func(io.Reader, []reflect.Type) ([]reflect.Value, error)

func (f DeserializerFunc) Deserialize(r io.Reader, t []reflect.Type) ([]reflect.Value, error) {
	return f(r, t)
}

// The UnmarshalerFunc type is an adapter to allow the use of ordinary
// functions as json.Unmarshaler[s]. If f is a function with the appropriate
// signature, UnmarshalerFunc(f) is a Unmarshaler that calls f.
type UnmarshalerFunc func([]byte) error

func (f UnmarshalerFunc) UnmarshalJSON(b []byte) error {
	return f(b)
}

// Unmarshaler returns a json.Unmarshaler which will unmarshal input data as the
// reflect type provided into the passed interface.
func Unmarshaler(pv *interface{}, t reflect.Type) json.Unmarshaler {
	return UnmarshalerFunc(func(b []byte) error {
		var v interface{}
		switch t.Kind() {
		case reflect.Ptr:
			v = reflect.New(t.Elem()).Interface()
		default:
			v = reflect.New(t).Interface()
		}
		if err := json.Unmarshal(b, v); err != nil {
			return err
		}
		if t.Kind() == reflect.Ptr {
			*pv = v
			return nil
		}
		*pv = reflect.ValueOf(v).Elem().Interface()
		return nil
	})
}

// JSONDeserializer deserializes data from the io.Reader as the types provided.
func JSONDeserializer(unmarshaler func(*interface{}, reflect.Type) json.Unmarshaler) Deserializer {
	return DeserializerFunc(func(r io.Reader, ts []reflect.Type) ([]reflect.Value, error) {
		var data []json.RawMessage
		if err := json.NewDecoder(r).Decode(&data); err != nil {
			return nil, err
		}
		if len(data) != len(ts) {
			return nil, errors.New("input data length doesn't match count of types provided")
		}
		vs := make([]reflect.Value, 0, len(ts))
		for i, t := range ts {
			var v interface{}
			if err := json.Unmarshal(data[i], unmarshaler(&v, t)); err != nil {
				return nil, err
			}
			vs = append(vs, reflect.ValueOf(v))
		}
		return vs, nil
	})
}

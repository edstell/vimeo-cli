package method

import (
	"encoding/json"
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
	return DeserializerFunc(func(r io.Reader, argt []reflect.Type) ([]reflect.Value, error) {
		var argr []json.RawMessage
		if err := json.NewDecoder(r).Decode(&argr); err != nil {
			return nil, err
		}
		// Append types to account for variadic arguments.
		itc := len(argr) - len(argt)
		for i := 0; i < itc; i++ {
			argt = append(argt, argt[len(argt)-1])
		}
		argv := make([]reflect.Value, 0, len(argr))
		for i, arg := range argr {
			var v interface{}
			if err := unmarshaler(&v, argt[i]).UnmarshalJSON(arg); err != nil {
				return nil, err
			}
			argv = append(argv, reflect.ValueOf(v))
		}
		return argv, nil
	})
}

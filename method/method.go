package method

import (
	"errors"
	"fmt"
	"io"
	"reflect"
)

// Caller wraps calling methods on a reflect.Value. By default it will serialize
// and deserialize data to-from json.
// User can modify the serializers via CallOption[s] when constructing the
// Caller.
type Caller struct {
	v reflect.Value
	s Serializer
	d Deserializer
}

// CallOption[s] are used to configure a Caller.
type CallOption func(*Caller)

// NewCaller initializes a Caller, applying any CallOption[s].
func NewCaller(v reflect.Value, opts ...CallOption) *Caller {
	caller := &Caller{
		v: v,
		s: JSONSerializer(Marshaler),
		d: JSONDeserializer(Unmarshaler),
	}
	for _, opt := range opts {
		opt(caller)
	}
	return caller
}

// Call the named method, reading input data from the io.Reader and writing
// output data to the io.Writer.
func (c *Caller) Call(name string, in io.Reader, out io.Writer) error {
	m, mt := c.v.MethodByName(name), c.v.Type()
	if !m.IsValid() {
		return errors.New(fmt.Sprintf("no method '%s' available", name))
	}
	ts := make([]reflect.Type, 0, mt.NumIn()-1)
	for i := 1; i < cap(ts); i++ {
		ts = append(ts, mt.In(i))
	}
	// TODO Additionally cater for variadic arguments.
	args, err := c.d.Deserialize(in, ts)
	if err != nil {
		return err
	}
	res := c.v.Call(args)
	return c.s.Serialize(out, res)
}
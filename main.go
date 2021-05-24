package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"

	"github.com/edstell/vimeo-cli/method"
	"github.com/silentsokolov/go-vimeo/vimeo"
	"golang.org/x/oauth2"
)

func exit(err error) {
	os.Stderr.WriteString(err.Error() + "\n")
	os.Exit(1)
}

func main() {
	client := vimeo.NewClient(oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: os.Getenv("VIMEO_ACCESS_TOKEN")},
	)), nil)
	args := os.Args[1:]
	if len(args) < 1 {
		exit(errors.New(fmt.Sprintf("vimeo [%s] method", strings.Join(serviceNames(client), " "))))
	}
	service := serviceByName(client, args[0])
	if !service.IsValid() {
		exit(errors.New(fmt.Sprintf("no method by name '%s'", args[0])))
	}
	if len(args) < 2 {
		exit(errors.New(fmt.Sprintf("vimeo %s [%s]", args[0], strings.Join(methodNames(service), " "))))
	}
	serializer := vimeoSerializer(method.JSONSerializer(method.Marshaler))
	caller := method.NewCaller(service, method.UsingSerializer(serializer))
	if err := caller.Call(args[1], os.Stdin, os.Stdout); err != nil {
		exit(err)
	}
}

func serviceNames(client *vimeo.Client) []string {
	v := reflect.ValueOf(client)
	t := reflect.TypeOf(client)
	names := make([]string, 0, v.Elem().NumField())
	for i := 0; i < t.Elem().NumField(); i++ {
		name := t.Elem().Field(i).Name
		field := v.Elem().FieldByName(name)
		if field.Kind() != reflect.Ptr {
			continue
		}
		if !field.Type().Elem().ConvertibleTo(reflect.TypeOf(vimeo.UsersService{})) {
			continue
		}
		names = append(names, name)
	}
	return names
}

func serviceByName(client *vimeo.Client, name string) reflect.Value {
	field := reflect.ValueOf(client).Elem().FieldByName(name)
	if !field.IsValid() {
		return reflect.Value{}
	}
	if field.Kind() != reflect.Ptr {
		return reflect.Value{}
	}
	if !field.Type().Elem().ConvertibleTo(reflect.TypeOf(vimeo.UsersService{})) {
		return reflect.Value{}
	}
	return field
}

func methodNames(v reflect.Value) []string {
	t := v.Type()
	names := make([]string, 0, t.NumMethod())
	for i := 0; i < t.NumMethod(); i++ {
		names = append(names, t.Method(i).Name)
	}
	return names
}

func vimeoSerializer(s method.Serializer) method.Serializer {
	errIface := reflect.TypeOf((*error)(nil)).Elem()
	return method.SerializerFunc(func(w io.Writer, res []reflect.Value) error {
		i := 0
		// Propagate or remove errors.
		for _, v := range res {
			ok := v.Type().Implements(errIface)
			if !ok {
				res[i] = v
				i++
				continue
			}
			if !v.IsNil() {
				return v.Interface().(error)
			}
		}
		res = res[:i]
		i = 0
		// Remove *vimeo.Response values.
		for _, v := range res {
			if v.Type() == reflect.TypeOf(&vimeo.Response{}) {
				continue
			}
			res[i] = v
			i++
		}
		res = res[:i]
		return s.Serialize(w, res)
	})
}

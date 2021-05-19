package vimeo

import (
	"fmt"
	"reflect"

	"github.com/silentsokolov/go-vimeo/vimeo"
)

var serviceType = reflect.TypeOf(vimeo.UsersService{})
var clientType = reflect.TypeOf(vimeo.Client{})

// Client provides a reflection wrapper around a vimeo client, to reference
// services in the client and call methods with strings passed from the CLI.
type Client struct {
	client reflect.Value
}

type Service struct {
	name    string
	service reflect.Value
}

// NewClient stores the vimeo.Client in our reflection wrapper.
func NewClient(client *vimeo.Client) *Client {
	return &Client{reflect.ValueOf(client).Elem()}
}

// Services returns the list of services available on the vimeo.Client.
func (c *Client) Services() []*Service {
	services := make([]*Service, 0, c.client.NumField())
	for i := 0; i < clientType.NumField(); i++ {
		fieldName := clientType.Field(i).Name
		fieldPtr := c.client.FieldByName(fieldName)
		if fieldPtr.Kind() != reflect.Ptr {
			continue
		}
		field := fieldPtr.Elem()
		if field.Kind() != reflect.Struct {
			continue
		}
		if !field.Type().ConvertibleTo(serviceType) {
			continue
		}
		services = append(services, &Service{
			name:    fieldName,
			service: field,
		})
	}
	return services
}

// Service looks up the named field, returning it in a Service reflection
// wrapper if the named field is found and is of type vimeo.service, otherwise
// nil is returned.
func (c *Client) Service(name string) *Service {
	fieldPtr := c.client.FieldByName(name)
	if !fieldPtr.IsValid() {
		return nil
	}
	if fieldPtr.Kind() != reflect.Ptr {
		return nil
	}
	fmt.Println("count methods: ", fieldPtr.NumMethod())
	field := fieldPtr.Elem()
	if !field.Type().ConvertibleTo(serviceType) {
		return nil
	}
	return &Service{
		name:    name,
		service: field,
	}
}

func (s *Service) Methods() []string {
	return nil
}

func (s *Service) String() string {
	return s.name
}

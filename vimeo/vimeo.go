package vimeo

import (
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

type Method struct {
	name   string
	method reflect.Value
}

// NewClient stores the vimeo.Client in our reflection wrapper.
func NewClient(client *vimeo.Client) *Client {
	return &Client{reflect.ValueOf(client)}
}

// Services returns the list of services available on the vimeo.Client.
func (c *Client) Services() []*Service {
	services := make([]*Service, 0, c.client.Elem().NumField())
	for i := 0; i < clientType.NumField(); i++ {
		fieldName := clientType.Field(i).Name
		field := c.client.Elem().FieldByName(fieldName)
		if field.Kind() != reflect.Ptr {
			continue
		}
		if !field.Type().Elem().ConvertibleTo(serviceType) {
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
	field := c.client.Elem().FieldByName(name)
	if !field.IsValid() {
		return nil
	}
	if field.Kind() != reflect.Ptr {
		return nil
	}
	if !field.Type().Elem().ConvertibleTo(serviceType) {
		return nil
	}
	return &Service{
		name:    name,
		service: field,
	}
}

// Methods lists the methods available for the given service.
func (s *Service) Methods() []reflect.Method {
	serviceType := s.service.Type()
	methods := make([]reflect.Method, 0, serviceType.NumMethod())
	for i := 0; i < serviceType.NumMethod(); i++ {
		methods = append(methods, serviceType.Method(i))
	}
	return methods
}

func (s *Service) String() string {
	return s.name
}

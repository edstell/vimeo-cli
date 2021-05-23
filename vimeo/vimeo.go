package vimeo

import (
	"errors"
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
func (c *Client) Service(name string) (*Service, error) {
	field := c.client.Elem().FieldByName(name)
	if !field.IsValid() {
		return nil, errors.New(fmt.Sprintf("'%s' is not a field in the vimeo client", name))
	}
	if field.Kind() != reflect.Ptr {
		return nil, errors.New("'%s' is not a field of type '*vimeo.service'")
	}
	if !field.Type().Elem().ConvertibleTo(serviceType) {
		return nil, errors.New("'%s' is not a field of type '*vimeo.service'")
	}
	return &Service{
		name:    name,
		service: field,
	}, nil
}

func (c *Client) Value() reflect.Value {
	return c.client
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

// Method returns the Service.Method for the given name.
func (s *Service) Method(name string) (reflect.Method, bool) {
	s.service.MethodByName(name)
	return s.service.Type().MethodByName(name)
}

func (s *Service) Value() reflect.Value {
	return s.service
}

func (s *Service) Name() string {
	return s.name
}

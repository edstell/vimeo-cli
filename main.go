package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"reflect"

	"github.com/edstell/morestrings"
	"github.com/edstell/vimeo-cli/vimeo"
	vimeoapi "github.com/silentsokolov/go-vimeo/vimeo"
	"golang.org/x/oauth2"
)

type Config struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	AccessToken  string `json:"access_token"`
}

const usage = "vimeo service methods [arguments...]"

func exit(err error) {
	os.Stderr.WriteString(err.Error() + "\n")
	os.Exit(1)
}

func readConfig(path string) (*Config, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	config := &Config{}
	if err := json.NewDecoder(f).Decode(config); err != nil {
		return nil, err
	}
	return config, nil
}

func main() {
	config, err := readConfig("config.json")
	if err != nil {
		exit(err)
	}
	client := vimeo.NewClient(vimeoapi.NewClient(oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.AccessToken},
	)), nil))
	if len(os.Args) < 2 {
		os.Stderr.WriteString("Available services:\n")
		services := client.Services()
		os.Stderr.WriteString(morestrings.JoinSlice(services, func(i int) string {
			return services[i].Name()
		}, "\n") + "\n")
		os.Exit(1)
	}
	service := client.Service(os.Args[1])
	if service == nil {
		exit(errors.New(usage))
	}
	if len(os.Args) < 3 {
		os.Stderr.WriteString("\nMethods for '" + service.Name() + "':\n")
		methods := service.Methods()
		os.Stderr.WriteString(morestrings.JoinSlice(methods, func(i int) string {
			return methods[i].Name
		}, "\n") + "\n")
		os.Exit(1)
	}
	method, ok := service.Method(os.Args[2])
	if !ok {
		exit(errors.New(usage))
	}
	args := []json.RawMessage{}
	if err := json.NewDecoder(os.Stdin).Decode(&args); err != nil {
		exit(err)
	}
	in := make([]reflect.Value, 0, method.Type.NumIn())
	in = append(in, service.Value())
	in = append(in, reflect.ValueOf(""))
	// for i := 1; i < cap(in) && i < len(args)-1; i++ {
	// 	v := reflect.Zero(method.Type.In(i)).Interface()
	// 	if err := json.Unmarshal(args[i-1], &v); err != nil {
	// 		exit(err)
	// 	}
	// 	in = append(in, reflect.ValueOf(v))
	// }
	res := method.Func.Call(in)
	if len(res) == 0 {
		os.Exit(0)
	}
	for _, v := range res {
		if err, ok := v.Interface().(error); ok && err != nil {
			exit(err)
		}
	}
	for _, v := range res {
		if err := json.NewEncoder(os.Stdout).Encode(v.Interface()); err != nil {
			// Skip.
		}
	}
}

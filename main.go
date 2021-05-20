package main

import (
	"context"
	"encoding/json"
	"errors"
	"os"

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
	args := os.Args[1:]
	if len(args) < 1 {
		os.Stderr.WriteString("Available services:\n")
		for _, service := range client.Services() {
			os.Stderr.WriteString(service.String() + "\n")
		}
		exit(errors.New(usage))
	}
	service := client.Service(args[0])
	if service == nil {
		exit(errors.New(usage))
	}
	if len(args) < 2 {
		os.Stderr.WriteString(usage + "\nMethods for '" + service.String() + "':\n")
		for _, method := range service.Methods() {
			os.Stderr.WriteString(method.Name + "\n")
		}
	}
	// for _, service := range client.Services() {
	// 	for _, method := range service.Methods() {
	// 		args := "("
	// 		for i := 1; i < method.Type.NumIn(); i++ {
	// 			if i != 1 {
	// 				args = args + ", "
	// 			}
	// 			args = args + method.Type.In(i).String()
	// 		}
	// 		fmt.Printf("%s.%s%s)\n", service, method.Name, args)
	// 	}
	// }
}

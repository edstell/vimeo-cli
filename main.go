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

func usage() error {
	return errors.New("vimeo service operation arguments...")
}

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
	// if len(os.Args[1:]) < 1 {
	// 	exit(usage())
	// }
	client.Service("Users")
}

package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"os"

	"github.com/silentsokolov/go-vimeo/vimeo"
	"golang.org/x/oauth2"
)

type Config struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	AccessToken  string `json:"access_token"`
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
	pathPtr := flag.String("config", "config.json", "path to client config")
	flag.Parse()
	fmt.Println(*pathPtr)
	config, err := readConfig(*pathPtr)
	if err != nil {
		exit(err)
	}
	client := vimeo.NewClient(oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(
		&oauth2.Token{AccessToken: config.AccessToken},
	)), nil)
	me, _, err := client.Users.Get("")
	if err != nil {
		exit(err)
	}
	fmt.Println(me)
}

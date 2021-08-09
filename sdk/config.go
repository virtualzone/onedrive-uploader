package sdk

import (
	"io/ioutil"
	"strings"
)

type Config struct {
	ClientID     string   `json:"client_id"`
	ClientSecret string   `json:"client_secret"`
	Scopes       []string `json:"scopes"`
	RedirectURL  string   `json:"redirect_uri"`
	SecretStore  string   `json:"secret_store"`
	Root         string   `json:"root"`
}

func ReadConfig(filename string) (*Config, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	var config Config
	if err := UnmarshalJSON(&config, data); err != nil {
		return nil, err
	}
	config.Root = strings.TrimSuffix(config.Root, "/")
	if !strings.HasPrefix(config.Root, "/") {
		config.Root = "/" + config.Root
	}
	return &config, nil
}

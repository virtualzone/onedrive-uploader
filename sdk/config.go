package sdk

import (
	"encoding/json"
	"os"
	"strings"
	"time"
)

type Config struct {
	ConfigFilePath string    `json:"-"`
	ClientID       string    `json:"client_id"`
	ClientSecret   string    `json:"client_secret"`
	Scopes         []string  `json:"scopes"`
	RedirectURL    string    `json:"redirect_uri"`
	Root           string    `json:"root"`
	AccessToken    string    `json:"access_token"`
	RefreshToken   string    `json:"refresh_token"`
	Expiry         time.Time `json:"expiry"`
	SecretStore    string    `json:"secret_store,omitempty"`
}

func ReadConfigData(data []byte) (*Config, error) {
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

func ReadConfig(filename string) (*Config, error) {
	data, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	config, err := ReadConfigData(data)
	config.ConfigFilePath = filename
	return config, err
}

func (config *Config) Write() error {
	data, err := json.Marshal(config)
	if err != nil {
		return err
	}
	if err := os.WriteFile(config.ConfigFilePath, data, 0600); err != nil {
		return err
	}
	return nil
}

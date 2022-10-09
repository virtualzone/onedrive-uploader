package main

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/virtualzone/onedrive-uploader/sdk"
)

type InteractiveConfig struct {
	TargetPath string
}

func (c *InteractiveConfig) readChar() rune {
	reader := bufio.NewReader(os.Stdin)
	char, _, err := reader.ReadRune()
	if err != nil {
		fmt.Printf("error reading from input: " + err.Error())
		os.Exit(1)
	}
	return char
}

func (c *InteractiveConfig) readString() string {
	reader := bufio.NewReader(os.Stdin)
	s, _ := reader.ReadString('\n')
	s = strings.TrimSuffix(s, "\n")
	s = strings.TrimSuffix(s, "\r")
	return s
}

func (c *InteractiveConfig) promptClientId(config *sdk.Config) {
	for config.ClientID == "" {
		fmt.Printf("OneDrive Client ID: ")
		config.ClientID = c.readString()
	}
}

func (c *InteractiveConfig) promptClientSecret(config *sdk.Config) {
	for config.ClientSecret == "" {
		fmt.Printf("OneDrive Client Secret: ")
		config.ClientSecret = c.readString()
	}
}

func (c *InteractiveConfig) promptScopes(config *sdk.Config) {
	fmt.Println("Available scopes:")
	fmt.Println("1) Default (Files.Read, Files.ReadWrite, Files.Read.All, Files.ReadWrite.All, offline_access)")
	fmt.Println("2) App Root (Files.ReadWrite.AppFolder, offline_access)")
	fmt.Println("3) Custom")
	scope := ' '
	for scope == ' ' {
		fmt.Printf("Select Scopes [1]: ")
		scope = c.readChar()
		switch scope {
		case '1', '\n':
			config.Scopes = []string{
				"Files.Read",
				"Files.ReadWrite",
				"Files.Read.All",
				"Files.ReadWrite.All",
				"offline_access",
			}
		case '2':
			config.Scopes = []string{
				"Files.ReadWrite.AppFolder",
				"offline_access",
			}
		case '3':
			fmt.Printf("Scopes: ")
			splitFunc := func(r rune) bool {
				return r == ',' || r == ' '
			}
			scopes := c.readString()
			scopesArray := strings.FieldsFunc(scopes, splitFunc)
			config.Scopes = make([]string, 0)
			for _, scopeItem := range scopesArray {
				scopeItem = strings.TrimSpace(scopeItem)
				if scopeItem != "" {
					config.Scopes = append(config.Scopes, scopeItem)
				}
			}
		default:
			scope = ' '
		}
	}
}

func (c *InteractiveConfig) promptRoot(config *sdk.Config) {
	fmt.Println("Available Drive Roots:")
	fmt.Println("1) Default (/drive/root)")
	fmt.Println("2) App Root (/drive/special/approot)")
	fmt.Println("3) Custom")
	root := ' '
	for root == ' ' {
		fmt.Printf("Select Drive Root [1]: ")
		root = c.readChar()
		switch root {
		case '1', '\n':
			config.Root = "/drive/root"
		case '2':
			config.Root = "/drive/special/approot"
		case '3':
			fmt.Printf("Custom Drive Root: ")
			config.Root = c.readString()
		default:
			root = ' '
		}
	}
}

func (c *InteractiveConfig) promptRedirectURL(config *sdk.Config) {
	fmt.Printf("Redirect URL? [http://localhost:53682/] ")
	config.RedirectURL = c.readString()
	if config.RedirectURL == "" {
		config.RedirectURL = "http://localhost:53682/"
	}
}

func (c *InteractiveConfig) promptSave(config *sdk.Config) {
	save := ""
	for save == "" {
		fmt.Printf("Save config? [" + c.TargetPath + "] ")
		save = c.readString()
		if save == "" {
			save = c.TargetPath
		}
		config.ConfigFilePath = save
		if err := config.Write(); err != nil {
			logError("Could not write config: " + err.Error())
			return
		}
		fmt.Printf("Config written to: %s\n", save)
	}
}

func (c *InteractiveConfig) Run() {
	config := &sdk.Config{}
	c.promptClientId(config)
	c.promptClientSecret(config)
	c.promptScopes(config)
	c.promptRoot(config)
	c.promptRedirectURL(config)
	c.promptSave(config)
}

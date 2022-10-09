package sdk

import (
	"fmt"
	"os"
	"runtime/debug"
	"testing"
)

const EnableIntegrationTests = true

var (
	RunIntegrationTests bool    = false
	IntegrationConfig   *Config = nil
	IntegrationClient   *Client = nil
)

func TestMain(m *testing.M) {
	if os.Getenv("INTEGRATION") == "1" || EnableIntegrationTests {
		fmt.Println("Integration tests enabled")
		RunIntegrationTests = true
		c, err := ReadConfig("../config.json")
		if err != nil {
			fmt.Println("Could not read config: " + err.Error())
			os.Exit(-1)
			return
		}
		IntegrationConfig = c
		client := CreateClient(c)
		if client.ShouldRenewAccessToken() {
			if _, err := client.RenewAccessToken(); err != nil {
				fmt.Println("Could not renew access token: " + err.Error())
				os.Exit(-1)
				return
			}
		}
		IntegrationClient = client
	}
	code := m.Run()
	os.Exit(code)
}

func checkTestBool(t *testing.T, expected, actual bool) {
	if expected != actual {
		t.Fatalf("Expected '%t', but got '%t' at:\n%s", expected, actual, debug.Stack())
	}
}

func checkTestString(t *testing.T, expected, actual string) {
	if expected != actual {
		t.Fatalf("Expected '%s', but got '%s' at:\n%s", expected, actual, debug.Stack())
	}
}

func checkTestInt(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Fatalf("Expected '%d', but got '%d' at:\n%s", expected, actual, debug.Stack())
	}
}

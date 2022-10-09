//go:build windows
// +build windows

package main

import (
	"os"
)

func GetConfigDir() string {
	return os.Getenv("APPDATA")
}

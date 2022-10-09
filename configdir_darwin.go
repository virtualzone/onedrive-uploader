//go:build darwin
// +build darwin

package main

import (
	"os"
)

func GetConfigDir() string {
	return os.Getenv("HOME") + "/Library/Application Support"
}

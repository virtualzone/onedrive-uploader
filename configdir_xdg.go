//go:build !windows && !darwin
// +build !windows,!darwin

package main

import (
	"os"
	"path/filepath"
)

func GetConfigDir() string {
	if os.Getenv("XDG_CONFIG_HOME") != "" {
		return os.Getenv("XDG_CONFIG_HOME")
	} else {
		return filepath.Join(os.Getenv("HOME"), ".config")
	}
}

package main

import (
	"errors"
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/virtualzone/onedrive-uploader/sdk"
)

type Flags struct {
	ConfigPath string
	Verbose    bool
	Quiet      bool
}

var (
	AppFlags = Flags{}
)

func printHelp() {
	flag.Usage()
	print("  config                             create config")
	print("  login                              perform login")
	print("  mkdir path                         create remote directory <path>")
	print("  ls path                            list items in <path>")
	print("  rm path                            delete <path>")
	print("  upload localFile path              upload <localFile> to <path>")
	print("  download sourceFile localPath      download <sourceFile> to <localPath>")
	print("  info path                          show info about <path>")
	print("  sha1 path                          get SHA1 hash for <path>")
	print("  sha256 path                        get SHA256 hash for <path>")
	print("  help                               show help")
	print("  migrate configPath                 migrate from old (< v0.6) config at <configPath> to current config")
	print("  version                            show version")
}

func prepareFlags() {
	flag.StringVar(&AppFlags.ConfigPath, "c", "", "path to config.json")
	flag.BoolVar(&AppFlags.Quiet, "q", false, "output errors only")
	flag.BoolVar(&AppFlags.Verbose, "v", false, "verbose output")
	flag.Parse()
}

func logVerbose(s string) {
	if AppFlags.Verbose {
		fmt.Println(s)
	}
}

func print(s string) {
	fmt.Println(s)
}

func log(s string) {
	if !AppFlags.Quiet {
		fmt.Println(s)
	}
}

func logError(s string) {
	fmt.Println(s)
	os.Exit(1)
}

func findConfigFilePath() (string, error) {
	if AppFlags.ConfigPath != "" {
		return AppFlags.ConfigPath, nil
	}
	configPath := GetConfigDir()
	if configPath == "" {
		var err error
		configPath, err = os.Getwd()
		if (err != nil) || (configPath == "") {
			return "", errors.New("could neither get system config dir nor current working dir")
		}
	} else {
		configPath = filepath.Join(configPath, "onedrive-uploader")
		_, err := os.Stat(configPath)
		if err != nil && os.IsNotExist(err) {
			os.MkdirAll(configPath, os.FileMode(0700))
		}
	}
	configPath = filepath.Join(configPath, "config.json")
	return configPath, nil
}

func main() {
	prepareFlags()
	cmd := ""
	if flag.NArg() > 0 {
		cmd = strings.ToLower(flag.Args()[0])
	}
	cmdDef := commands[cmd]
	if cmdDef == nil {
		print("OneDrive Uploader " + AppVersion)
		printHelp()
		return
	}
	logVerbose("OneDrive Uploader " + AppVersion)
	args := []string{}
	if flag.NArg() > 1 {
		args = flag.Args()[1:]
	}
	if len(args) < cmdDef.MinArgs {
		printHelp()
		return
	}
	outputRenderer := &OutputRenderer{
		Quiet: AppFlags.Quiet,
	}
	var client *sdk.Client = nil
	if cmdDef.RequireConfig {
		configPath, err := findConfigFilePath()
		if (err != nil) || (configPath == "") {
			logError("Could not initialize config path: " + err.Error())
			return
		}
		conf, err := sdk.ReadConfig(configPath)
		if err != nil {
			logError("Could not read config: " + err.Error())
			return
		}
		client = sdk.CreateClient(conf)
		client.UseTransferSignals = true
		client.Verbose = AppFlags.Verbose
		if cmdDef.InitSecretStore {
			logVerbose("Reading secret store...")
			if client.ShouldRenewAccessToken() {
				logVerbose("Renewing access token...")
				outputRenderer.initSpinner("Renewing access token...")
				if _, err := client.RenewAccessToken(); err != nil {
					outputRenderer.stopSpinner()
					logError("Could not renew access token: " + err.Error())
					return
				}
				outputRenderer.stopSpinner()
			}
		}
	}
	cmdDef.Fn(client, outputRenderer, args)
}

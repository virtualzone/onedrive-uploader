package main

import (
	"flag"
	"fmt"
	"os"
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
	print("  version                            show version")
}

func prepareFlags() {
	flag.StringVar(&AppFlags.ConfigPath, "c", "./config.json", "path to config.json")
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
		conf, err := sdk.ReadConfig(AppFlags.ConfigPath)
		if err != nil {
			logError("Could not read config: " + err.Error())
			return
		}
		client = sdk.CreateClient(conf)
		client.UseTransferSignals = true
		client.Verbose = AppFlags.Verbose
		if cmdDef.InitSecretStore {
			logVerbose("Reading secret store...")
			if err := client.ReadSecretStore(); err != nil {
				logError("Could not read secret store: " + err.Error())
				return
			}
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

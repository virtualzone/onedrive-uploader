package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/virtualzone/onedrive-uploader/sdk"
)

type Flags struct {
	Verbose bool
	Quiet   bool
}

var (
	AppFlags = Flags{}
)

func printHelp() {
	flag.Usage()
	log("  login                              perform login")
	log("  mkdir [path]                       create remote directory <path>")
	log("  ls [path]                          list items in <path>")
	log("  rm [path]                          delete <path>")
	log("  upload [localFile] [path]          upload <localFile> to <path>")
	log("  download [sourceFile] [localPath]  download <sourceFile> to <localPath>")
	log("  help                               show help")
	log("  version                            show version")
}

func prepareFlags() {
	flag.BoolVar(&AppFlags.Quiet, "q", false, "output errors only")
	flag.BoolVar(&AppFlags.Verbose, "v", false, "verbose output")
	flag.Parse()
}

func logVerbose(s string) {
	if AppFlags.Verbose {
		fmt.Println(s)
	}
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
		log("OneDrive Uploader " + AppVersion)
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
	conf, err := sdk.ReadConfig("./config.json")
	if err != nil {
		logError("Could not read config: " + err.Error())
		return
	}
	client := sdk.CreateClient(conf)
	if cmdDef.InitSecretStore {
		logVerbose("Reading secret store...")
		if err := client.ReadSecretStore(); err != nil {
			logError("Could not read secret store: " + err.Error())
			return
		}
		if client.ShouldRenewAccessToken() {
			logVerbose("Renewing access token...")
			if _, err := client.RenewAccessToken(); err != nil {
				logError("Could not renew access token: " + err.Error())
				return
			}
		}
	}
	cmdDef.Fn(client, args)
}

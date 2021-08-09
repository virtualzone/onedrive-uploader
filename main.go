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
	log("  login                        Perform login")
	log("  mkdir [path]                 Create remote directory <path>")
	log("  upload [path] [localFile]    Upload <localFile> to <path>")
	log("  rm [path]                    Delete <path>")
	log("  help                         Show commands")
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
	os.Exit(2)
}

func main() {
	prepareFlags()
	logVerbose("OneDrive Uploader " + AppVersion)
	cmd := ""
	if flag.NArg() > 0 {
		cmd = strings.ToLower(flag.Args()[0])
	}
	cmdDef := commands[cmd]
	if cmdDef == nil {
		printHelp()
		return
	}
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
		logVerbose("Renewing access token...")
		if _, err := client.RenewAccessToken(); err != nil {
			logError("Could not renew access token: " + err.Error())
			return
		}
	}
	cmdDef.Fn(client, args)
}

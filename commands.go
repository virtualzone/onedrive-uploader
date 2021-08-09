package main

import "github.com/virtualzone/onedrive-uploader/sdk"

type CommandFunction func(client *sdk.Client, args []string)

type CommandFunctionDefinition struct {
	Fn              CommandFunction
	MinArgs         int
	InitSecretStore bool
}

var (
	commands = map[string]*CommandFunctionDefinition{
		"login":  {Fn: cmdLogin, MinArgs: 0, InitSecretStore: false},
		"mkdir":  {Fn: cmdCreateDir, MinArgs: 1, InitSecretStore: true},
		"upload": {Fn: cmdUpload, MinArgs: 2, InitSecretStore: true},
		"rm":     {Fn: cmdDelete, MinArgs: 1, InitSecretStore: true},
	}
)

func cmdLogin(client *sdk.Client, args []string) {
	log("------------------------------------")
	log("Open a browser and go to:")
	log(client.GetLoginURL())
	log("------------------------------------")
	log("Waiting for code...")
	if err := client.Login(); err != nil {
		logError("Could not log in: " + err.Error())
		return
	}
	log("Login successful.")
}

func cmdCreateDir(client *sdk.Client, args []string) {
	if err := client.CreateDir(args[0]); err != nil {
		logError("Could not create folder: " + err.Error())
		return
	}
	log("Folder created.")
}

func cmdUpload(client *sdk.Client, args []string) {
	if err := client.Upload(args[0], args[1]); err != nil {
		logError("Could not upload file: " + err.Error())
		return
	}
	log("File uploaded.")
}

func cmdDelete(client *sdk.Client, args []string) {
	if err := client.Delete(args[0]); err != nil {
		logError("Could not delete: " + err.Error())
		return
	}
	log("Deleted.")
}

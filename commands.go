package main

import (
	"strconv"

	"github.com/virtualzone/onedrive-uploader/sdk"
)

type CommandFunction func(client *sdk.Client, args []string)

type CommandFunctionDefinition struct {
	Fn              CommandFunction
	MinArgs         int
	InitSecretStore bool
}

var (
	commands = map[string]*CommandFunctionDefinition{
		"login":    {Fn: cmdLogin, MinArgs: 0, InitSecretStore: false},
		"mkdir":    {Fn: cmdCreateDir, MinArgs: 1, InitSecretStore: true},
		"upload":   {Fn: cmdUpload, MinArgs: 2, InitSecretStore: true},
		"download": {Fn: cmdDownload, MinArgs: 2, InitSecretStore: true},
		"rm":       {Fn: cmdDelete, MinArgs: 1, InitSecretStore: true},
		"ls":       {Fn: cmdList, MinArgs: 1, InitSecretStore: true},
		"info":     {Fn: cmdInfo, MinArgs: 1, InitSecretStore: true},
		"sha1":     {Fn: cmdSHA1, MinArgs: 1, InitSecretStore: true},
		"sha256":   {Fn: cmdSHA256, MinArgs: 1, InitSecretStore: true},
		"version":  {Fn: cmdVersion, MinArgs: 0, InitSecretStore: false},
	}
)

func cmdLogin(client *sdk.Client, args []string) {
	log("------------------------------------")
	log("Open a browser and go to:")
	print(client.GetLoginURL())
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

func cmdDownload(client *sdk.Client, args []string) {
	if err := client.Download(args[0], args[1]); err != nil {
		logError("Could not download file: " + err.Error())
		return
	}
	log("File downloaded.")
}

func cmdDelete(client *sdk.Client, args []string) {
	if err := client.Delete(args[0]); err != nil {
		logError("Could not delete: " + err.Error())
		return
	}
	log("Deleted.")
}

func cmdList(client *sdk.Client, args []string) {
	list, err := client.List(args[0])
	if err != nil {
		logError("Could not list: " + err.Error())
		return
	}
	for _, item := range list {
		itemType := "d"
		if item.File.MimeType != "" {
			itemType = "f"
		}
		print(itemType + " " + item.Name)
	}
}

func cmdInfo(client *sdk.Client, args []string) {
	item, err := client.Info(args[0])
	if err != nil {
		logError("Could not get info: " + err.Error())
		return
	}
	itemType := "folder"
	if item.File.MimeType != "" {
		itemType = "file"
	}
	print("Type:           " + itemType)
	print("Size:           " + strconv.FormatInt(item.SizeBytes, 10) + " bytes")
	if itemType == "folder" {
		print("Child Count:    " + strconv.Itoa(item.Folder.ChildCount))
	} else {
		print("MIME Type:      " + item.File.MimeType)
		print("SHA1 Hash:      " + item.File.Hashes.SHA1)
		print("SHA256 Hash:    " + item.File.Hashes.SHA256)
		print("Quick XOR Hash: " + item.File.Hashes.QuickXOR)
	}
}

func cmdSHA1(client *sdk.Client, args []string) {
	item, err := client.Info(args[0])
	if err != nil {
		logError("Could not get info: " + err.Error())
		return
	}
	print(item.File.Hashes.SHA1)
}

func cmdSHA256(client *sdk.Client, args []string) {
	item, err := client.Info(args[0])
	if err != nil {
		logError("Could not get info: " + err.Error())
		return
	}
	print(item.File.Hashes.SHA256)
}

func cmdVersion(client *sdk.Client, args []string) {
	print(AppVersion)
}

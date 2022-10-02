package main

import (
	"io/fs"
	"strconv"
	"time"

	"github.com/virtualzone/onedrive-uploader/sdk"
)

type CommandFunction func(client *sdk.Client, renderer *OutputRenderer, args []string)

type CommandFunctionDefinition struct {
	Fn              CommandFunction
	MinArgs         int
	InitSecretStore bool
	RequireConfig   bool
}

var (
	commands = map[string]*CommandFunctionDefinition{
		"login":    {Fn: cmdLogin, MinArgs: 0, InitSecretStore: false, RequireConfig: true},
		"mkdir":    {Fn: cmdCreateDir, MinArgs: 1, InitSecretStore: true, RequireConfig: true},
		"upload":   {Fn: cmdUpload, MinArgs: 2, InitSecretStore: true, RequireConfig: true},
		"download": {Fn: cmdDownload, MinArgs: 2, InitSecretStore: true, RequireConfig: true},
		"rm":       {Fn: cmdDelete, MinArgs: 1, InitSecretStore: true, RequireConfig: true},
		"ls":       {Fn: cmdList, MinArgs: 1, InitSecretStore: true, RequireConfig: true},
		"info":     {Fn: cmdInfo, MinArgs: 1, InitSecretStore: true, RequireConfig: true},
		"sha1":     {Fn: cmdSHA1, MinArgs: 1, InitSecretStore: true, RequireConfig: true},
		"sha256":   {Fn: cmdSHA256, MinArgs: 1, InitSecretStore: true, RequireConfig: true},
		"version":  {Fn: cmdVersion, MinArgs: 0, InitSecretStore: false, RequireConfig: false},
	}
)

func cmdLogin(client *sdk.Client, renderer *OutputRenderer, args []string) {
	log("------------------------------------")
	log("Open a browser and go to:")
	print(client.GetLoginURL())
	log("------------------------------------")
	renderer.initSpinner("Waiting for code...")
	err := client.Login()
	renderer.stopSpinner()
	if err != nil {
		logError("Could not log in: " + err.Error())
		return
	}
	log("Login successful.")
}

func cmdCreateDir(client *sdk.Client, renderer *OutputRenderer, args []string) {
	renderer.initSpinner("Creating directory...")
	err := client.CreateDir(args[0])
	renderer.stopSpinner()
	if err != nil {
		logError("Could not create folder: " + err.Error())
		return
	}
	log("Folder created.")
}

func cmdUpload(client *sdk.Client, renderer *OutputRenderer, args []string) {
	done := false
	go func() {
		fileStat := <-client.ChannelTransferStart
		renderer.initProgressBar(fileStat.Size(), "Uploading "+fileStat.Name()+"...")
	}()
	go func() {
		for !done {
			bytes := <-client.ChannelTransferProgress
			renderer.updateProgressBar(bytes)
			time.Sleep(50 * time.Millisecond)
		}
	}()
	go func() {
		done = <-client.ChannelTransferFinish
	}()
	err := client.Upload(args[0], args[1])
	if err != nil {
		logError("Could not upload file: " + err.Error())
		return
	}
	log("File uploaded.")
}

func cmdDownload(client *sdk.Client, renderer *OutputRenderer, args []string) {
	done := false
	go func() {
		var fileStat fs.FileInfo = nil
		for fileStat == nil {
			fileStat := <-client.ChannelTransferStart
			if fileStat == nil {
				renderer.initSpinner("Retrieving information...")
			} else {
				renderer.stopSpinner()
				renderer.initProgressBar(fileStat.Size(), "Downloading "+fileStat.Name()+"...")
			}
		}
	}()
	go func() {
		for !done {
			bytes := <-client.ChannelTransferProgress
			renderer.updateProgressBar(bytes)
			time.Sleep(50 * time.Millisecond)
		}
	}()
	go func() {
		done = <-client.ChannelTransferFinish
	}()
	err := client.Download(args[0], args[1])
	if err != nil {
		logError("Could not download file: " + err.Error())
		return
	}
	log("File downloaded.")
}

func cmdDelete(client *sdk.Client, renderer *OutputRenderer, args []string) {
	renderer.initSpinner("Deleting...")
	err := client.Delete(args[0])
	renderer.stopSpinner()
	if err != nil {
		logError("Could not delete: " + err.Error())
		return
	}
	log("Deleted.")
}

func cmdList(client *sdk.Client, renderer *OutputRenderer, args []string) {
	renderer.initSpinner("Retrieving directory listing...")
	list, err := client.List(args[0])
	renderer.stopSpinner()
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

func cmdInfo(client *sdk.Client, renderer *OutputRenderer, args []string) {
	renderer.initSpinner("Retrieving information...")
	item, err := client.Info(args[0])
	renderer.stopSpinner()
	if err != nil {
		logError("Could not get info: " + err.Error())
		return
	}
	itemType := "folder"
	if item.Type == sdk.DriveItemTypeFile {
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

func cmdSHA1(client *sdk.Client, renderer *OutputRenderer, args []string) {
	renderer.initSpinner("Retrieving SHA1 hash...")
	item, err := client.Info(args[0])
	renderer.stopSpinner()
	if err != nil {
		logError("Could not get info: " + err.Error())
		return
	}
	print(item.File.Hashes.SHA1)
}

func cmdSHA256(client *sdk.Client, renderer *OutputRenderer, args []string) {
	renderer.initSpinner("Retrieving SHA256 hash...")
	item, err := client.Info(args[0])
	renderer.stopSpinner()
	if err != nil {
		logError("Could not get info: " + err.Error())
		return
	}
	print(item.File.Hashes.SHA256)
}

func cmdVersion(client *sdk.Client, renderer *OutputRenderer, args []string) {
	print(AppVersion)
}

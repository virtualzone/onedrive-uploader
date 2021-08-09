package main

import (
	"log"
	"os"
	"strings"

	"github.com/virtualzone/onedrive-uploader/sdk"
)

func printHelp() {
	log.Println("Usage:")
	log.Println("    login                        Perform login")
	log.Println("    mkdir [path]                 Create remote directory <path>")
	log.Println("    upload [path] [localFile]    Upload <localFile> to <path>")
	log.Println("    help                         Show commands")
}

func main() {
	log.Println("OneDrive Uploader")
	args := os.Args[1:]
	if len(args) == 0 {
		printHelp()
		return
	}
	log.Println("Reading config...")
	conf, err := sdk.ReadConfig("./config.json")
	if err != nil {
		log.Fatalln("Could not read config")
	}
	client := sdk.CreateClient(conf)
	cmd := strings.ToLower(args[0])
	if cmd == "login" {
		log.Println("------------------------------------")
		log.Println("Open a browser and go to:")
		log.Println(client.GetLoginURL())
		log.Println("------------------------------------")
		log.Println("Waiting for code...")
		client.Login()
	}
	if cmd == "help" {
		printHelp()
	}
	if cmd == "mkdir" {
		if len(args) < 2 {
			printHelp()
			return
		}
		client.ReadSecretStore()
		log.Println("Renewing access token...")
		client.RenewAccessToken()
		client.CreateDir(args[1])
	}
	if cmd == "upload" {
		if len(args) < 3 {
			printHelp()
			return
		}
		client.ReadSecretStore()
		log.Println("Renewing access token...")
		client.RenewAccessToken()
		client.Upload(args[1], args[2])
	}
}

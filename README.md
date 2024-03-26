# OneDrive Uploader
[![](https://img.shields.io/github/v/release/virtualzone/onedrive-uploader)](https://github.com/virtualzone/onedrive-uploader/releases)
[![](https://img.shields.io/github/release-date/virtualzone/onedrive-uploader)](https://github.com/virtualzone/onedrive-uploader/releases)
[![](https://img.shields.io/github/actions/workflow/status/virtualzone/onedrive-uploader/test.yml?branch=main)](https://github.com/virtualzone/onedrive-uploader/actions)
[![Go Report Card](https://goreportcard.com/badge/github.com/virtualzone/onedrive-uploader)](https://goreportcard.com/report/github.com/virtualzone/onedrive-uploader)
[![](https://img.shields.io/github/license/virtualzone/onedrive-uploader)](https://github.com/virtualzone/onedrive-uploader/blob/master/LICENSE)

Command line (CLI) utility for uploading files to Microsoft OneDrive using the Microsoft Graph REST API.

## Features
* Upload, download and delete files
* Create and delete directories
* List folder contents
* Get information (including SHA1 and SHA256 hashes) for drive items
* Supports "special folders" (such as App Folder / App Root)
* Pre-compiled binaries on Linux, MacOS and Windows

![](https://raw.githubusercontent.com/virtualzone/onedrive-uploader/main/res/progress.gif)

## Getting started

### 1. Create an Azure Application
OneDrive Uploader requires an application to be registered with Microsoft. This application is not exposed anywhere. It just serves as a means to authenticate access to your OneDrive with specified permissions.

1. Log in to the [Microsoft Azure Portal](https://portal.azure.com/).
1. Navigate to "App registrations".
1. Create a new application with supported account type "Accounts in any organizational directory (Any Azure AD directory - Multitenant) and personal Microsoft accounts (e.g. Skype, Xbox)" and the following Web redirect URL: http://localhost:53682/
1. Copy the Application (client) ID.
1. Navigate to "Certificates & secrets", create a new Client secret and copy the Secret Value (*not* the ID).
1. Navigate to "API permissions", click "Add permission", choose "Microsoft Graph", select "Delegated". Then search and add the required permissions:
    1. Access to App Folder only: ```Files.ReadWrite.AppFolder, offline_access, User.Read```
    1. Access to entire OneDrive: ```Files.Read, Files.ReadWrite, Files.Read.All, Files.ReadWrite.All, offline_access, User.Read```

### 2. Download the binary
Download the appropriate binary for your operating system and architecture from the [Releases page](https://github.com/virtualzone/onedrive-uploader/releases) and make it executable.

You can use the following one-liner to perform the necessary steps on Linux and MacOS:
```
curl -s -L https://git.io/JRie0 | bash
```

On MacOS with [Homebrew](https://brew.sh) installed:
```
brew install virtualzone/tap/onedrive-uploader
```

If you want to build the binaries from source instead, clone the repository and execute ```make```. This requires Go to be installed.
```
git clone https://github.com/virtualzone/onedrive-uploader.git
cd onedrive-uploader
make
```

### 3. Configuration & Login 
Run the following command to create the configuration file:
```
onedrive-uploader config
```

After that, execute the following command to log in with your OneDrive account:
```
onedrive-uploader login
```

For headless machines you must perform the actual login on a computer *with* a web browser. To do this, you can...
* ...either run the ```config``` and ```login``` commands on another computer with a web browser and then copy the ```config.json``` to the headless computer after having logged in
* ...or forward port 53682 from your computer with a web brower to your headless machine, e.g. by using SSH: ```ssh -L 53682:headless_ip:53682 user@headless_ip```
* ...or use the ```curl``` command with fallback url

The configuration file is stored in the following directory (if not specified otherwise using the ```-c``` parameter):

* Linux: ```${HOME}/.config/onedrive-uploader```
* MacOS: ```${HOME}/Library/Application Support/onedrive-uploader```
* Windows: ```${APPDATA}/onedrive-uploader```

## Commands and example usage
Create a new remote directory named "test":
```
onedrive-uploader mkdir test
```

Create a new remote directory named "test2" below the "test" folder:
```
onedrive-uploader mkdir test/test2
```

List contents of the "test" folder:
```
onedrive-uploader ls test
```

Upload local file "image.jpg" to the "test" folder:
```
onedrive-uploader upload /tmp/image.jpg test
```

Download "notes.docx" from the root directory:
```
onedrive-uploader download /notes.docx /tmp
```

Delete "notes.docx" from the root directory:
```
onedrive-uploader rm /notes.docx
```

Get information about file "notes.docx" in folder "test":
```
onedrive-uploader info /test/notes.docx
```

Get SHA1 hash for file "notes.docx" in folder "test":
```
onedrive-uploader sha1 /test/notes.docx
```

Get SHA256 hash for file "notes.docx" in folder "test":
```
onedrive-uploader sha256 /test/notes.docx
```

Print help and available commands:
```
onedrive-uploader help
```

Use config file at a specific path:
```
onedrive-uploader -c /path/to/config.json mkdir test
```

### Important note for users of version < 0.6
The configuration file format and path has changed as of version 0.6.

Please run the following command to migrate your existing configuration file to the new file format and location after having installed the latest version of OneDrive Uploader:

```
onedrive-uploader migrate /path/to/existing/config.json
```

You can safely delete your existing ```config.json``` and ```secret.json``` files after the migration has been performed successfully.

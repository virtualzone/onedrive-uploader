# OneDrive Uploader
[![](https://img.shields.io/github/v/release/virtualzone/onedrive-uploader)](https://github.com/virtualzone/onedrive-uploader/releases)
[![](https://img.shields.io/github/release-date/virtualzone/onedrive-uploader)](https://github.com/virtualzone/onedrive-uploader/releases)
[![](https://img.shields.io/github/workflow/status/virtualzone/onedrive-uploader/build-release)](https://github.com/virtualzone/onedrive-uploader/actions)
[![](https://img.shields.io/github/license/virtualzone/onedrive-uploader)](https://github.com/virtualzone/onedrive-uploader/blob/master/LICENSE)

Command line (CLI) utility for uploading files to Microsoft OneDrive using the Microsoft Graph REST API.

## Features
* Upload, download and delete files
* Create and delete directories
* List folder contents
* Get information (including SHA1 and SHA256 hashes) for drive items
* Supports "special folders" (such as App Folder / App Root)
* Pre-compiled binaries on Linux, MacOS and Windows

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
Download the appropriate binary for your operating system and architecture from the [Releases page](https://github.com/virtualzone/onedrive-uploader/releases) and make it executable. You can use the following one-liner to perform the necessary steps:
```
curl -s -L https://git.io/JRie0 | bash
```

If you want to build the binaries from source instead, clone the repository and execute ```make```. This requires Go to be installed.
```
git clone https://github.com/virtualzone/onedrive-uploader.git
cd onedrive-uploader
make
```

### 3. Create a configuration file
In the binary's folder, create a configuration file named ```config.json```.

Example ```config.json``` for an App with access to a user's entire OneDrive:
```
{
    "client_id": "019ccb8b-118f-4559-ad2c-2ccda5b9def6",
    "client_secret": "some-client-secret",
    "scopes": [
        "Files.Read",
        "Files.ReadWrite",
        "Files.Read.All",
        "Files.ReadWrite.All",
        "offline_access"
    ],
    "redirect_uri": "http://localhost:53682/",
    "secret_store": "./secret.json",
    "root": "/drive/root"
}
```

Example ```config.json``` for an App registered with permissions to access the App's Folder:
```
{
    "client_id": "019ccb8b-118f-4559-ad2c-2ccda5b9def6",
    "client_secret": "some-client-secret",
    "scopes": [
        "Files.ReadWrite.AppFolder",
        "offline_access"
    ],
    "redirect_uri": "http://localhost:53682/",
    "secret_store": "./secret.json",
    "root": "/drive/special/approot"
}
```

You can set the following properties:
* ```client_id```: Application (client) ID from Azure App registration
* ```client_secret```: Client Secret Value from Azure App registration
* ```scopes```: Permission scopes
* ```redirect_uri```: Redirect URL (required for login only, must match the URL set in your Azure App)
* ```secret_store```: Path to a file where the access and refresh tokens will be stored
* ```root```: Root folder within your OneDrive

### 4. Perform login
To log in with your OneDrive account, execute the following command:
```
onedrive-uploader login
```

For headless machines you must perform the actual login on a computer *with* a web browser. To do this, you can...
* ...either have the same ```config.json``` on your headless machine and your computer with a web browser and then copy the ```secret.json``` to the headless computer after having logged in
* ...or forward port 53682 from your computer with a web brower to your headless machine, e.g. by using SSH: ```ssh -L 53682:headless_ip:53682 user@headless_ip```

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
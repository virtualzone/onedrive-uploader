# OneDrive Uploader
Command line (CLI) utility for uploading files to Microsoft OneDrive using the Microsoft Graph REST API.

## Features
* Upload, download and delete files
* Create and delete directories
* List folder contents
* Supports "special folders" (such as App Folder / App Root)
* Pre-compiled binaries on Linux, MacOS and Windows

## Getting started

### 1. Create an Azure Application
OneDrive Uploader requires an application to be registered with Microsoft. This application is not exposed anywhere. It just serves as a means to authenticate access to your OneDrive with specified permissions.

1. Log in to the [Microsoft Azure Portal](https://portal.azure.com/).
1. Navigate to "App registrations".
1. Create a new application with the following redirect URL: http://localhost:53682/
1. Copy the Application (client) ID.
1. Navigate to "Certificates & secrets", create a new Client secret and copy the Secret Value (*not* the ID).
1. Navigate to "API permissions", click "Add permission", choose "Microsoft Graph", select "Delegated". Then search and add the required permissions:
    1. Access to App Folder only: ```Files.ReadWrite.AppFolder, offline_access, User.Read```
    1. Access to entire OneDrive: ```Files.Read, Files.ReadWrite, Files.Read.All, Files.ReadWrite.All, offline_access, User.Read```

### 2. Download the binary
Download the appropriate binary for your operating system from the [Releases page](https://github.com/virtualzone/onedrive-uploader/releases) and make it executable.

Example:
```
curl https://github.com/virtualzone/onedrive-uploader/releases/download/v0.2.0/onedrive-uploader_linux_amd64_v0.2.0 --output onedrive-uploader
chmod +x onedrive-uploader
```

### 3. Create a configuration file
In the binary's folder, create a configuration file named ```config.json```.

The following is an example for an App registered with permissions to access the App's Folder:

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
To log in with your OneDrive account, execute the following command on a computer which has a web browser running. On headless machines, perform this task on a computer *with* a web browser and copy the ```secret.json``` the headless computer afterwards.

```
onedrive-uploader login
```

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

Print help and available commands:
```
onedrive-uploader help
```

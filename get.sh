#!/bin/bash
if ! command -v curl &> /dev/null
then
    echo "curl is required to install OneDrive Uploader, but could not be found."
    exit 1
fi
SUDO=""
if [ "$(id -u)" != "0" ]; then
    SUDO="sudo"
    if ! command -v $SUDO &> /dev/null
    then
        SUDO="doas"
    fi
    if ! command -v $SUDO &> /dev/null
    then
        echo "Neither sudo nor doas found, try running installer script as root."
        exit 1
    fi
fi
OS=`uname -s | tr '[A-Z]' '[a-z]'`
if [[ "$OS" == "darwin" ]]; then
    OS="macos"
fi
ARCH=`uname -m | tr '[A-Z]' '[a-z]'`
if [[ "$ARCH" == "aarch64" || "$ARCH" == "aarch64_be" || "$ARCH" == "armv8b" || "$ARCH" == "armv8l"  ]]; then
    ARCH="arm64"
elif [[ "$ARCH" == "x86_64" ]]; then
    ARCH="amd64"
elif [[ "$ARCH" == "armv71" ]]; then
    ARCH="arm"
fi
URL=`curl -s https://api.github.com/repos/virtualzone/onedrive-uploader/releases/latest | grep "browser_download_url" | grep "_${OS}_${ARCH}_" | cut -d : -f 2,3 | tr -d \" | xargs`
if [[ "$URL" == "" ]]; then
    echo "Could not find binary for OS '$OS' and architecture '$ARCH'."
    echo "Please check for an appropriate binary at: https://github.com/virtualzone/onedrive-uploader"
    exit 1
fi
if [ "$(id -u)" != "0" ]; then
    echo "Please specify your sudo password when asked. It's required to write the binary to: /usr/local/bin/"
fi
$SUDO curl -s -L "${URL}" -o /usr/local/bin/onedrive-uploader && \
    $SUDO chmod +x /usr/local/bin/onedrive-uploader &&
    VERSION=`/usr/local/bin/onedrive-uploader version` &&
    echo "Successfully installed OneDrive Uploader ${VERSION}."
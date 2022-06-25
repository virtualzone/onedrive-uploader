VERSION=`cat ./VERSION | awk NF`

all: clean update_version_file linux macos windows

clean:
	rm -f build/*

update_version_file:
	echo "package main\n\nvar AppVersion = \"${VERSION}\"" > version.go

linux: linux_amd64 linux_arm64 linux_arm

macos: macos_amd64 macos_arm64

windows: windows_amd64

linux_amd64:
	env CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o build/onedrive-uploader_linux_amd64_${VERSION}

linux_arm64:
	env CGO_ENABLED=0 GOOS=linux GOARCH=arm64 go build -ldflags="-w -s" -o build/onedrive-uploader_linux_arm64_${VERSION}

linux_arm:
	env CGO_ENABLED=0 GOOS=linux GOARCH=arm go build -ldflags="-w -s" -o build/onedrive-uploader_linux_arm_${VERSION}

macos_amd64:
	env CGO_ENABLED=0 GOOS=darwin GOARCH=amd64 go build -ldflags="-w -s" -o build/onedrive-uploader_macos_amd64_${VERSION}

macos_arm64:
	env CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="-w -s" -o build/onedrive-uploader_macos_arm64_${VERSION}

windows_amd64:
	env CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-w -s" -o build/onedrive-uploader_windows_amd64_${VERSION}
VERSION=0.0.1

all: clean linux macos windows

clean:
	rm -f build/*

linux: linux_amd64 linux_arm64 linux_arm

macos: macos_amd64 macos_arm64

windows: windows_amd64

linux_amd64:
	env GOOS=linux GOARCH=amd64 go build -o build/onedrive-uploader_linux_amd64_${VERSION}

linux_arm64:
	env GOOS=linux GOARCH=arm64 go build -o build/onedrive-uploader_linux_arm64_${VERSION}

linux_arm:
	env GOOS=linux GOARCH=arm go build -o build/onedrive-uploader_linux_arm_${VERSION}

macos_amd64:
	env GOOS=darwin GOARCH=amd64 go build -o build/onedrive-uploader_macos_amd64_${VERSION}

macos_arm64:
	env GOOS=darwin GOARCH=arm64 go build -o build/onedrive-uploader_macos_arm64_${VERSION}

windows_amd64:
	env GOOS=windows GOARCH=amd64 go build -o build/onedrive-uploader_windows_amd64_${VERSION}
package sdk

import (
	"errors"
	"io/ioutil"
	"mime"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

var (
	UploadSessionFileSizeLimit int = 4 * 1000 * 1000
	UploadSessionMultiple      int = 320 * 1024
	UploadSessionRangeSize     int = UploadSessionMultiple * 10
)

type UploadSessionResponse struct {
	UploadURL string    `json:"uploadUrl"`
	Expiry    time.Time `json:"expirationDateTime"`
}

type EmptyStruct struct{}

func (client *Client) Upload(targetFolder, localFilePath string) error {
	fileName := filepath.Base(localFilePath)
	if fileName == "" || fileName == "." || fileName == ".." {
		return errors.New("please specify a file, not a directory")
	}
	targetFolder = strings.TrimPrefix(strings.TrimSuffix(targetFolder, "/"), "/")
	if !strings.HasSuffix(targetFolder, "/") {
		targetFolder += "/"
	}
	if !strings.HasPrefix(targetFolder, "/") {
		targetFolder = "/" + targetFolder
	}
	fileStat, err := os.Stat(localFilePath)
	if err != nil {
		return err
	}
	mimeType := mime.TypeByExtension(filepath.Ext(localFilePath))
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}
	if fileStat.Size() < int64(UploadSessionFileSizeLimit) {
		// Use simple upload
		return client.uploadSimple(fileName, mimeType, targetFolder, localFilePath)
	}
	// Use upload session
	session, err := client.startUploadSession(fileName, targetFolder)
	if err != nil {
		return err
	}
	return client.uploadToSession(session.UploadURL, mimeType, localFilePath, fileStat.Size())
}

func (client *Client) uploadToSession(uploadUrl, mimeType, localFilePath string, fileSize int64) error {
	data := make([]byte, UploadSessionRangeSize)
	f, err := os.Open(localFilePath)
	if err != nil {
		return err
	}
	defer f.Close()
	var offset int64 = 0
	n := 0
	for offset < fileSize {
		n, _ = f.ReadAt(data, offset)
		if n < UploadSessionRangeSize {
			data = append([]byte(nil), data[:n]...)
		}
		status, _, err := client.httpSendFilePart("PUT", uploadUrl, mimeType, offset, int64(n), fileSize, data)
		if err != nil {
			return err
		}
		if status != http.StatusAccepted && status != http.StatusCreated {
			return errors.New("received unexpected status code " + strconv.Itoa(status))
		}
		offset += int64(n)
	}
	return nil
}

func (client *Client) startUploadSession(fileName, targetFolder string) (*UploadSessionResponse, error) {
	//log.Printf("Creating upload session for '%s' to '%s'...\n", fileName, targetFolder)
	url := GraphURL + "me" + client.Config.Root + ":" + targetFolder + fileName + ":/createUploadSession"
	payload := &EmptyStruct{}
	status, data, err := client.httpPostJSON(url, payload)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, errors.New("received unexpected status code " + strconv.Itoa(status))
	}
	var uploadSession UploadSessionResponse
	if err := UnmarshalJSON(&uploadSession, data); err != nil {
		return nil, err
	}
	return &uploadSession, nil
}

func (client *Client) uploadSimple(fileName, mimeType, targetFolder, localFilePath string) error {
	data, err := ioutil.ReadFile(localFilePath)
	if err != nil {
		return err
	}
	//log.Printf("Uploading '%s' (%s) to '%s'...\n", fileName, mimeType, targetFolder)
	url := GraphURL + "me" + client.Config.Root + ":" + targetFolder + fileName + ":/content"
	status, _, err := client.httpSendFile("PUT", url, mimeType, data)
	if err != nil {
		return err
	}
	if status != http.StatusCreated {
		return errors.New("received unexpected status code " + strconv.Itoa(status))
	}
	return nil
}

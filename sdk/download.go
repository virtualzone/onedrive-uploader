package sdk

import (
	"errors"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type DownloadFileStat struct {
	FileName  string
	SizeBytes int64
}

func (s *DownloadFileStat) Name() string {
	return s.FileName
}

func (s *DownloadFileStat) Size() int64 {
	return s.SizeBytes
}

func (s *DownloadFileStat) Mode() fs.FileMode {
	return 0
}

func (s *DownloadFileStat) ModTime() time.Time {
	return time.Now()
}

func (s *DownloadFileStat) IsDir() bool {
	return false
}

func (s *DownloadFileStat) Sys() any {
	return nil
}

func (client *Client) Download(sourceFilePath, targetFolder string) error {
	if len(sourceFilePath) > 0 && sourceFilePath[0] == '.' {
		return errors.New("invalid source path (should start with /)")
	}
	fileName := filepath.Base(sourceFilePath)
	if fileName == "" || fileName == "." || fileName == ".." {
		return errors.New("please specify a file, not a directory")
	}
	sourceFilePath = strings.TrimPrefix(sourceFilePath, "/")
	if !strings.HasPrefix(sourceFilePath, "/") {
		sourceFilePath = "/" + sourceFilePath
	}
	if !strings.HasSuffix(targetFolder, "/") {
		targetFolder += "/"
	}
	// Get file info
	client.signalTransferStart(nil)
	info, err := client.Info(sourceFilePath)
	if err != nil {
		return err
	}
	// Start download
	url := GraphURL + "me" + client.Config.Root + ":" + sourceFilePath + ":/content"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+client.Config.AccessToken)
	httpClient := &http.Client{}
	resp, err := httpClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusNotFound {
		return errors.New("file not found")
	}
	if resp.StatusCode != http.StatusOK {
		data, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		return client.handleResponseError(resp.StatusCode, data)
	}
	fileStat := &DownloadFileStat{
		FileName:  info.Name,
		SizeBytes: info.SizeBytes,
	}
	client.signalTransferStart(fileStat)
	out, err := os.Create(targetFolder + fileName)
	if err != nil {
		client.signalTransferFinish()
		return err
	}
	defer out.Close()
	total := int64(0)
	reader := &ProgressReader{
		Reader: resp.Body,
		OnReadProgress: func(r int64) {
			total += r
			client.signalTransferProgress(total)
		},
	}
	_, err = io.Copy(out, reader)
	client.signalTransferFinish()
	if err != nil {
		return err
	}
	return nil
}

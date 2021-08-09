package sdk

import (
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func (client *Client) Download(sourceFilePath, targetFolder string) error {
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
	url := GraphURL + "me" + client.Config.Root + ":" + sourceFilePath + ":/content"
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return err
	}
	req.Header.Add("Authorization", "Bearer "+client.SecretStore.AccessToken)
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
		return errors.New("received unexpected status code " + strconv.Itoa(resp.StatusCode))
	}
	out, err := os.Create(targetFolder + fileName)
	if err != nil {
		return err
	}
	defer out.Close()
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		return err
	}
	return nil
}

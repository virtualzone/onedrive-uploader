package sdk

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
)

func (client *Client) Delete(path string) error {
	path = strings.TrimSuffix(path, "/")
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	url := GraphURL + "me" + client.Config.Root + ":" + path
	status, err := client.httpDelete(url)
	if err != nil {
		return err
	}
	if status == http.StatusNotFound {
		return errors.New("path not found")
	}
	if status != http.StatusNoContent {
		return errors.New("received unexpected status code " + strconv.Itoa(status))
	}
	return nil
}

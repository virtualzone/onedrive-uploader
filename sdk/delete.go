package sdk

import (
	"errors"
	"net/http"
	"strings"
)

func (client *Client) Delete(path string) error {
	if len(path) > 0 && path[0] == '.' {
		return errors.New("invalid path (should start with /)")
	}
	path = strings.TrimSuffix(path, "/")
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	url := GraphURL + "me" + client.Config.Root + ":" + path
	status, data, err := client.httpDelete(url)
	if err != nil {
		return err
	}
	if status == http.StatusNotFound {
		return errors.New("path not found")
	}
	if status != http.StatusNoContent {
		return client.handleResponseError(status, data)
	}
	return nil
}

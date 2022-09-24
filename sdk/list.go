package sdk

import (
	"errors"
	"net/http"
	"strings"
)

type ListResponse struct {
	Items    []DriveItem `json:"value"`
	NextLink string      `json:"@odata.nextLink"`
}

func (client *Client) List(path string) ([]*DriveItem, error) {
	if len(path) > 0 && path[0] == '.' {
		return nil, errors.New("invalid path (should start with /)")
	}
	path = strings.TrimSuffix(path, "/")
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	url := GraphURL + "me" + client.Config.Root + ":" + path + ":/children"
	if path == "/" {
		url = GraphURL + "me" + client.Config.Root + "/children"
	}
	params := map[string]string{
		"top":     "100000",
		"orderby": "name",
	}
	status, data, err := client.httpGet(url, params)
	if err != nil {
		return nil, err
	}
	if status == http.StatusNotFound {
		return nil, errors.New("path not found")
	}
	if status != http.StatusOK {
		return nil, client.handleResponseError(status, data)
	}
	var resp ListResponse
	if err := UnmarshalJSON(&resp, data); err != nil {
		return nil, err
	}
	var result []*DriveItem
	for i := range resp.Items {
		driveItem := &resp.Items[i]
		if driveItem.File.MimeType != "" {
			driveItem.Type = DriveItemTypeFile
		} else {
			driveItem.Type = DriveItemTypeFolder
		}
		result = append(result, driveItem)
	}
	return result, nil
}

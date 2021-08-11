package sdk

import (
	"errors"
	"net/http"
	"strconv"
	"strings"
)

type ListResponse struct {
	Items    []DriveItem `json:"value"`
	NextLink string      `json:"@odata.nextLink"`
}

func (client *Client) List(path string) ([]*DriveItem, error) {
	path = strings.TrimSuffix(path, "/")
	if !strings.HasPrefix(path, "/") {
		path = "/" + path
	}
	url := GraphURL + "me" + client.Config.Root + ":" + path + ":/children"
	if path == "/" {
		url = GraphURL + "me" + client.Config.Root + "/children"
	}
	params := map[string]string{
		"$top":     "100000",
		"$orderby": "name",
	}
	status, data, err := client.httpGet(url, params)
	if err != nil {
		return nil, err
	}
	if status == http.StatusNotFound {
		return nil, errors.New("path not found")
	}
	if status != http.StatusOK {
		return nil, errors.New("received unexpected status code " + strconv.Itoa(status))
	}
	var resp ListResponse
	if err := UnmarshalJSON(&resp, data); err != nil {
		return nil, err
	}
	var result []*DriveItem
	for i := range resp.Items {
		result = append(result, &resp.Items[i])
	}
	return result, nil
}

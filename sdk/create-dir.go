package sdk

import (
	"errors"
	"net/http"
	"strings"
)

type FolderProperties struct {
}

type CreateFolderRequest struct {
	Name             string           `json:"name"`
	Folder           FolderProperties `json:"folder"`
	ConflictBehavior string           `json:"@microsoft.graph.conflictBehavior"`
}

func (client *Client) CreateDir(path string) error {
	if len(path) > 0 && path[0] == '.' {
		return errors.New("invalid path (should start with /)")
	}
	path = strings.TrimPrefix(strings.TrimSuffix(path, "/"), "/")
	pathParts := strings.Split(path, "/")
	newFolder := pathParts[len(pathParts)-1]
	parentPath := "/"
	if len(pathParts) > 1 {
		parentPathParts := pathParts[:len(pathParts)-1]
		parentPath = "/" + strings.Join(parentPathParts, "/")
	}
	req := &CreateFolderRequest{
		Name:             strings.TrimSpace(newFolder),
		Folder:           FolderProperties{},
		ConflictBehavior: "fail",
	}
	url := GraphURL + "me" + client.Config.Root + ":" + parentPath + ":/children"
	if parentPath == "/" {
		url = GraphURL + "me" + client.Config.Root + "/children"
	}
	status, data, err := client.httpPostJSON(url, req)
	if err != nil {
		return err
	}
	if status == http.StatusConflict {
		// Already exists
		return nil
	}
	if status != http.StatusCreated {
		return client.handleResponseError(status, data)
	}
	return nil
}

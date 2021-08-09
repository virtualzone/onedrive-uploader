package sdk

import (
	"errors"
	"net/http"
	"strconv"
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
	path = strings.TrimPrefix(strings.TrimSuffix(path, "/"), "/")
	pathParts := strings.Split(path, "/")
	newFolder := pathParts[len(pathParts)-1]
	parentPath := "/"
	if len(pathParts) > 1 {
		parentPathParts := pathParts[:len(pathParts)-1]
		parentPath = "/" + strings.Join(parentPathParts, "/")
	}
	//log.Printf("Creating folder '%s' in parent path '%s'...\n", newFolder, parentPath)
	req := &CreateFolderRequest{
		Name:             newFolder,
		Folder:           FolderProperties{},
		ConflictBehavior: "rename",
	}
	url := GraphURL + "me" + client.Config.Root + ":" + parentPath + ":/children"
	if parentPath == "/" {
		url = GraphURL + "me" + client.Config.Root + "/children"
	}
	status, _, err := client.httpPostJSON(url, req)
	if err != nil {
		return err
	}
	if status != http.StatusCreated {
		return errors.New("received unexpected status code " + strconv.Itoa(status))
	}
	return nil
}

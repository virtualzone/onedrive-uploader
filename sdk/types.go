package sdk

import "time"

type FileHashes struct {
	QuickXOR string `json:"quickXorHash"`
	SHA1     string `json:"sha1Hash"`
	SHA256   string `json:"sha256Hash"`
}

type FileItem struct {
	MimeType string     `json:"mimeType"`
	Hashes   FileHashes `json:"hashes"`
}

type FileSystemInfo struct {
	Created      time.Time `json:"createdDateTime"`
	LastModified time.Time `json:"lastModifiedDateTime"`
}

type FolderItem struct {
	ChildCount int `json:"childCount"`
}

type DriveItem struct {
	Name           string         `json:"name"`
	SizeBytes      int64          `json:"size"`
	File           FileItem       `json:"file"`
	Folder         FolderItem     `json:"folder"`
	FileSystemInfo FileSystemInfo `json:"fileSystemInfo"`
}

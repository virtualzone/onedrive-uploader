package sdk

import (
	"crypto/rand"
	"crypto/sha1"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"strings"
	"testing"

	"github.com/google/uuid"
)

func TestCreateDeleteDir(t *testing.T) {
	dirName := "/test-" + uuid.New().String()

	// Create
	err := IntegrationClient.CreateDir(dirName)
	checkTestBool(t, err == nil, true)

	// Delete
	err = IntegrationClient.Delete(dirName)
	checkTestBool(t, err == nil, true)
}

func TestDirInfo(t *testing.T) {
	dirName := "/test-" + uuid.New().String()
	testCreateDir(t, dirName, dirName)
}

func TestDirLeadingSpace(t *testing.T) {
	dirName := uuid.New().String()
	testCreateDir(t, "/ test-"+dirName, "/test-"+dirName)
}

func TestDirTrailingSpace(t *testing.T) {
	dirName := uuid.New().String()
	testCreateDir(t, "/test-"+dirName+" ", "/test-"+dirName)
}

func TestUploadDownloadSmall(t *testing.T) {
	fileName := uuid.New().String() + ".txt"
	testUpload(t, fileName, fileName, 1, "text/plain; charset=utf-8")
}

func TestUploadDownloadLarge(t *testing.T) {
	fileName := uuid.New().String() + ".dat"
	fileSizeKB := (UploadSessionFileSizeLimit / 1024 * 2) + 10
	testUpload(t, fileName, fileName, fileSizeKB, "application/octet-stream")
}

func TestList(t *testing.T) {
	// Create
	dirName := "/test-" + uuid.New().String()
	err := IntegrationClient.CreateDir(dirName)
	checkTestBool(t, err == nil, true)

	// Check empty dir
	items, err := IntegrationClient.List(dirName)
	checkTestBool(t, err == nil, true)
	checkTestInt(t, 0, len(items))

	// Prepare local files
	fileName1 := "a_" + uuid.New().String() + ".txt"
	fileName2 := "b_" + uuid.New().String() + ".png"
	createRandomFile("/tmp/"+fileName1, 1)
	createRandomFile("/tmp/"+fileName2, 2)
	defer os.Remove("/tmp/" + fileName1)
	defer os.Remove("/tmp/" + fileName2)

	// Create sub-folder and upload two files
	IntegrationClient.CreateDir(dirName + "/sub")
	IntegrationClient.Upload("/tmp/"+fileName1, dirName)
	IntegrationClient.Upload("/tmp/"+fileName2, dirName)

	// Check dir again
	items, err = IntegrationClient.List(dirName)
	checkTestBool(t, err == nil, true)
	checkTestInt(t, 3, len(items))
	// 1
	checkTestString(t, fileName1, items[0].Name)
	checkTestBool(t, true, items[0].Type == DriveItemTypeFile)
	checkTestInt(t, 1024, int(items[0].SizeBytes))
	// 2
	checkTestString(t, fileName2, items[1].Name)
	checkTestBool(t, true, items[1].Type == DriveItemTypeFile)
	checkTestInt(t, 2048, int(items[1].SizeBytes))
	// 3
	checkTestString(t, "sub", items[2].Name)
	checkTestBool(t, true, items[2].Type == DriveItemTypeFolder)

	// Delete one file
	IntegrationClient.Delete(dirName + "/" + fileName2)

	// Check dir again
	items, err = IntegrationClient.List(dirName)
	checkTestBool(t, err == nil, true)
	checkTestInt(t, 2, len(items))
	checkTestString(t, fileName1, items[0].Name)
	checkTestString(t, "sub", items[1].Name)

	// Delete
	err = IntegrationClient.Delete(dirName)
	checkTestBool(t, err == nil, true)
}

func TestUploadLeadingWhitespace(t *testing.T) {
	fileName := uuid.New().String() + ".txt"
	testUpload(t, " "+fileName, fileName, 1, "text/plain; charset=utf-8")
}

func TestUploadTrailingWhitespace(t *testing.T) {
	fileName := uuid.New().String() + ".txt"
	testUpload(t, fileName+" ", fileName, 1, "text/plain; charset=utf-8")
}

func TestUploadInvalidChars(t *testing.T) {
	const chars = "~\"#%&*:<>?\\{|}"
	for i := 0; i < len(chars); i++ {
		c := string(chars[i])
		fileName := uuid.New().String()
		testUpload(t, fileName+"_"+c+"_some.txt", fileName+"_"+"_"+"_some.txt", 1, "text/plain; charset=utf-8")
	}
}

func testUpload(t *testing.T, localFileName, expectedRemoteFileName string, sizeKB int, mimeType string) {
	dirName := "/test-" + uuid.New().String()
	err := createRandomFile("/tmp/"+localFileName, sizeKB)
	checkTestBool(t, true, err == nil)
	defer os.Remove("/tmp/" + localFileName)
	hash1, _ := getSHA1Hash("/tmp/" + localFileName)
	hash256, _ := getSHA256Hash("/tmp/" + localFileName)

	// Create
	err = IntegrationClient.CreateDir(dirName)
	checkTestBool(t, err == nil, true)

	// Upload
	err = IntegrationClient.Upload("/tmp/"+localFileName, dirName)
	checkTestBool(t, err == nil, true)

	// Get info on folder
	item, err := IntegrationClient.Info(dirName)
	checkTestBool(t, err == nil, true)
	checkTestBool(t, item.Type == DriveItemTypeFolder, true)
	checkTestString(t, strings.TrimPrefix(dirName, "/"), item.Name)
	checkTestInt(t, 1, item.Folder.ChildCount)

	// Get info on file
	item, err = IntegrationClient.Info(dirName + "/" + expectedRemoteFileName)
	checkTestBool(t, err == nil, true)
	checkTestBool(t, item.Type == DriveItemTypeFile, true)
	checkTestString(t, expectedRemoteFileName, item.Name)
	checkTestString(t, mimeType, item.File.MimeType)
	checkTestInt(t, sizeKB*1024, int(item.SizeBytes))
	checkTestString(t, strings.ToUpper(hash1), item.File.Hashes.SHA1)
	checkTestString(t, strings.ToUpper(hash256), item.File.Hashes.SHA256)

	// Download
	os.Remove("/tmp/" + localFileName)
	IntegrationClient.Download(dirName+"/"+expectedRemoteFileName, "/tmp")
	hash256Downloaded, _ := getSHA256Hash("/tmp/" + expectedRemoteFileName)
	checkTestString(t, hash256, hash256Downloaded)

	// Delete
	err = IntegrationClient.Delete(dirName)
	checkTestBool(t, err == nil, true)
}

func testCreateDir(t *testing.T, dirName string, expectedDirName string) {
	// Create
	err := IntegrationClient.CreateDir(dirName)
	checkTestBool(t, err == nil, true)

	// Get info
	item, err := IntegrationClient.Info(expectedDirName)
	checkTestBool(t, err == nil, true)
	checkTestBool(t, item.Type == DriveItemTypeFolder, true)
	checkTestString(t, strings.TrimPrefix(expectedDirName, "/"), item.Name)
	checkTestInt(t, 0, item.Folder.ChildCount)

	// Delete
	err = IntegrationClient.Delete(expectedDirName)
	checkTestBool(t, err == nil, true)
}

func createRandomFile(fileName string, sizeKB int) error {
	file, err := os.Create(fileName)
	if err != nil {
		return err
	}
	defer file.Close()
	data := make([]byte, 1024)
	for i := 1; i <= sizeKB; i++ {
		rand.Read(data)
		if _, err := file.Write(data); err != nil {
			return err
		}
	}
	return nil
}

func getSHA1Hash(fileName string) (string, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha1.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

func getSHA256Hash(fileName string) (string, error) {
	f, err := os.Open(fileName)
	if err != nil {
		return "", err
	}
	defer f.Close()
	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", err
	}
	return fmt.Sprintf("%x", h.Sum(nil)), nil
}

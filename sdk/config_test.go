package sdk

import "testing"

func TestReadConfigData(t *testing.T) {
	data := `{
		"client_id": "019ccb8b-118f-4559-ad2c-2ccda5b9def6",
		"client_secret": "some-client-secret",
		"scopes": [
			"Files.Read",
			"Files.ReadWrite",
			"Files.Read.All",
			"Files.ReadWrite.All",
			"offline_access"
		],
		"redirect_uri": "http://localhost:53682/",
		"secret_store": "./secret.json",
		"root": "drive/root/"
	}`
	c, err := ReadConfigData([]byte(data))
	checkTestBool(t, true, err == nil)
	checkTestString(t, "/drive/root", c.Root)
}

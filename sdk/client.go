package sdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var (
	GraphURL = "https://graph.microsoft.com/v1.0/"
)

type Client struct {
	Config      *Config
	SecretStore *SecretStore
}

type HTTPRequestParams map[string]string

func CreateClient(conf *Config) *Client {
	client := &Client{
		Config:      conf,
		SecretStore: &SecretStore{},
	}
	return client
}

func UnmarshalJSON(o interface{}, body []byte) error {
	if body == nil {
		return errors.New("body is NIL")
	}
	if err := json.Unmarshal(body, &o); err != nil {
		return err
	}
	return nil
}

func (client *Client) buildURIParams(params HTTPRequestParams) string {
	if params != nil {
		var b strings.Builder
		for name, value := range params {
			if b.Len() != 0 {
				fmt.Fprint(&b, "&")
			}
			fmt.Fprintf(&b, "%s=%s", name, url.QueryEscape(value))
		}
		return b.String()
	}
	return ""
}

func (client *Client) buildURI(uri string, params HTTPRequestParams) string {
	paramString := client.buildURIParams(params)
	if paramString != "" {
		uri += "?" + paramString
	}
	return uri
}

func (client *Client) httpPostForm(uri string, params HTTPRequestParams) (int, []byte, error) {
	requestHeaders := make(HTTPRequestParams)
	requestHeaders["Content-Type"] = "application/x-www-form-urlencoded"
	payload := []byte(client.buildURIParams(params))
	return client.httpRequest("POST", uri, requestHeaders, nil, payload)
}

func (client *Client) httpSendFile(method, uri, mimeType string, data []byte) (int, []byte, error) {
	requestHeaders := make(HTTPRequestParams)
	requestHeaders["Content-Type"] = mimeType
	if client.SecretStore.AccessToken != "" {
		requestHeaders["Authorization"] = "Bearer " + client.SecretStore.AccessToken
	}
	return client.httpRequest(method, uri, requestHeaders, nil, data)
}

func (client *Client) httpSendFilePart(method, uri, mimeType string, offset, n, fileSize int64, data []byte) (int, []byte, error) {
	requestHeaders := make(HTTPRequestParams)
	requestHeaders["Content-Type"] = mimeType
	requestHeaders["Content-Length"] = strconv.FormatInt(n, 10)
	requestHeaders["Content-Range"] = "bytes " + strconv.FormatInt(offset, 10) + "-" + strconv.FormatInt(n+offset-1, 10) + "/" + strconv.FormatInt(fileSize, 10)
	if client.SecretStore.AccessToken != "" {
		requestHeaders["Authorization"] = "Bearer " + client.SecretStore.AccessToken
	}
	return client.httpRequest(method, uri, requestHeaders, nil, data)
}

func (client *Client) httpSendJSON(method, uri string, o interface{}) (int, []byte, error) {
	payload, err := json.Marshal(o)
	if err != nil {
		return -1, nil, err
	}
	requestHeaders := make(HTTPRequestParams)
	requestHeaders["Content-Type"] = "application/json"
	if client.SecretStore.AccessToken != "" {
		requestHeaders["Authorization"] = "Bearer " + client.SecretStore.AccessToken
	}
	return client.httpRequest(method, uri, requestHeaders, nil, payload)
}

func (client *Client) httpDelete(uri string) (int, error) {
	requestHeaders := make(HTTPRequestParams)
	if client.SecretStore.AccessToken != "" {
		requestHeaders["Authorization"] = "Bearer " + client.SecretStore.AccessToken
	}
	status, _, err := client.httpRequest("DELETE", uri, requestHeaders, nil, nil)
	return status, err
}

func (client *Client) httpGet(uri string, params HTTPRequestParams) (int, []byte, error) {
	requestHeaders := make(HTTPRequestParams)
	if client.SecretStore.AccessToken != "" {
		requestHeaders["Authorization"] = "Bearer " + client.SecretStore.AccessToken
	}
	return client.httpRequest("GET", uri, requestHeaders, params, nil)
}

func (client *Client) httpPostJSON(uri string, o interface{}) (int, []byte, error) {
	return client.httpSendJSON("POST", uri, o)
}

func (client *Client) httpPutJSON(uri string, o interface{}) (int, []byte, error) {
	return client.httpSendJSON("PUT", uri, o)
}

func (client *Client) httpRequest(method, uri string, requestHeaders, params HTTPRequestParams, payload []byte) (int, []byte, error) {
	httpClient := &http.Client{}
	uri = client.buildURI(uri, params)
	//log.Println(method + " request for: " + uri)
	reader := bytes.NewReader(payload)
	req, err := http.NewRequest(method, uri, reader)
	if err != nil {
		return -1, nil, err
	}
	for name, value := range requestHeaders {
		req.Header.Add(name, value)
	}
	resp, err := httpClient.Do(req)
	if err != nil {
		return -1, nil, err
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return -1, nil, err
	}
	//log.Println(string(body))
	return resp.StatusCode, body, nil
}

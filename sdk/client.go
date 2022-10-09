package sdk

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var (
	GraphURL = "https://graph.microsoft.com/v1.0/"
)

type transferProgress func(int64)

type Client struct {
	Config                  *Config
	Verbose                 bool
	UseTransferSignals      bool
	ChannelTransferStart    chan fs.FileInfo
	ChannelTransferProgress chan int64
	ChannelTransferFinish   chan bool
}

type HTTPRequestParams map[string]string

func CreateClient(conf *Config) *Client {
	client := &Client{
		Config:             conf,
		Verbose:            false,
		UseTransferSignals: false,
	}
	client.ResetChannels()
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

func IsHTTPStatusOK(status int) bool {
	return (200 <= status && status <= 299)
}

func (client *Client) ResetChannels() {
	if client.ChannelTransferStart != nil {
		close(client.ChannelTransferStart)
	}
	if client.ChannelTransferProgress != nil {
		close(client.ChannelTransferProgress)
	}
	if client.ChannelTransferFinish != nil {
		close(client.ChannelTransferFinish)
	}
	client.ChannelTransferStart = make(chan fs.FileInfo)
	client.ChannelTransferProgress = make(chan int64)
	client.ChannelTransferFinish = make(chan bool)
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
	return client.httpRequest("POST", uri, requestHeaders, nil, payload, nil)
}

func (client *Client) httpSendFile(method, uri, mimeType string, data []byte, progress transferProgress) (int, []byte, error) {
	requestHeaders := make(HTTPRequestParams)
	requestHeaders["Content-Type"] = mimeType
	if client.Config.AccessToken != "" {
		requestHeaders["Authorization"] = "Bearer " + client.Config.AccessToken
	}
	return client.httpRequest(method, uri, requestHeaders, nil, data, progress)
}

func (client *Client) httpSendFilePart(method, uri, mimeType string, offset, n, fileSize int64, data []byte, progress transferProgress) (int, []byte, error) {
	requestHeaders := make(HTTPRequestParams)
	requestHeaders["Content-Type"] = mimeType
	requestHeaders["Content-Length"] = strconv.FormatInt(n, 10)
	requestHeaders["Content-Range"] = "bytes " + strconv.FormatInt(offset, 10) + "-" + strconv.FormatInt(n+offset-1, 10) + "/" + strconv.FormatInt(fileSize, 10)
	if client.Config.AccessToken != "" {
		requestHeaders["Authorization"] = "Bearer " + client.Config.AccessToken
	}
	return client.httpRequest(method, uri, requestHeaders, nil, data, progress)
}

func (client *Client) httpSendJSON(method, uri string, o interface{}) (int, []byte, error) {
	payload, err := json.Marshal(o)
	if err != nil {
		return -1, nil, err
	}
	requestHeaders := make(HTTPRequestParams)
	requestHeaders["Content-Type"] = "application/json"
	if client.Config.AccessToken != "" {
		requestHeaders["Authorization"] = "Bearer " + client.Config.AccessToken
	}
	return client.httpRequest(method, uri, requestHeaders, nil, payload, nil)
}

func (client *Client) httpDelete(uri string) (int, []byte, error) {
	requestHeaders := make(HTTPRequestParams)
	if client.Config.AccessToken != "" {
		requestHeaders["Authorization"] = "Bearer " + client.Config.AccessToken
	}
	return client.httpRequest("DELETE", uri, requestHeaders, nil, nil, nil)
}

func (client *Client) httpGet(uri string, params HTTPRequestParams) (int, []byte, error) {
	requestHeaders := make(HTTPRequestParams)
	if client.Config.AccessToken != "" {
		requestHeaders["Authorization"] = "Bearer " + client.Config.AccessToken
	}
	return client.httpRequest("GET", uri, requestHeaders, params, nil, nil)
}

func (client *Client) httpPostJSON(uri string, o interface{}) (int, []byte, error) {
	return client.httpSendJSON("POST", uri, o)
}

func (client *Client) httpRequest(method, uri string, requestHeaders, params HTTPRequestParams, payload []byte, progress transferProgress) (int, []byte, error) {
	httpClient := &http.Client{}
	uri = client.buildURI(uri, params)
	total := int64(0)
	reader := &ProgressReader{
		Reader: bytes.NewReader(payload),
		OnReadProgress: func(r int64) {
			total += r
			if progress != nil {
				progress(total)
			}
		},
	}
	req, err := http.NewRequest(method, uri, reader)
	req.ContentLength = int64(reader.Len())
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
	return resp.StatusCode, body, nil
}

func (client *Client) handleResponseError(status int, data []byte) error {
	var resp ErrorResponse
	if err := UnmarshalJSON(&resp, data); err != nil {
		return errors.New("received unexpected status code " + strconv.Itoa(status))
	}
	return errors.New("received unexpected status code " + strconv.Itoa(status) + ": " + resp.Error.Message + " (" + resp.Error.Code + ")")
}

func (client *Client) signalTransferStart(info fs.FileInfo) {
	if !client.UseTransferSignals {
		return
	}
	client.ChannelTransferStart <- info
}

func (client *Client) signalTransferProgress(b int64) {
	if !client.UseTransferSignals {
		return
	}
	client.ChannelTransferProgress <- b
}

func (client *Client) signalTransferFinish() {
	if !client.UseTransferSignals {
		return
	}
	client.ChannelTransferFinish <- true
}

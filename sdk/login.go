package sdk

import (
	"context"
	"errors"
	"net/http"
	"strings"
	"time"
)

type LoginRedeemCodeResponse struct {
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

const (
	loginHTMLResponseHeader = "<!doctype html>" +
		"<html>" +
		"<head>" +
		"<title>OneDrive Uploader</title>" +
		"</head>" +
		"<body>" +
		"<h1>OneDrive Uploader</h1>"
	loginHTMLResponseFooter = "</body></html>"
	loginHTMLResponseOK     = "<p>Received authorization code from Microsoft Graph API.</p>" +
		"<p>Please return to your terminal now.</p>"
	loginHTMLResponseNotFound = "<p>Error: Page not found.</p>"
)

func (client *Client) Login() error {
	code := client.expectCode()
	grant, err := client.redeemCodeForAccessToken(code)
	if err != nil {
		return err
	}
	if grant.AccessToken == "" {
		return errors.New("received empty access token")
	}
	return nil
}

func (client *Client) UpdateSecretStore(grant *LoginRedeemCodeResponse) error {
	expiry := time.Now().Add(time.Second * time.Duration(grant.ExpiresIn))
	client.Config.AccessToken = grant.AccessToken
	client.Config.RefreshToken = grant.RefreshToken
	client.Config.Expiry = expiry
	return client.Config.Write()
}

func (client *Client) GetLoginURL() string {
	params := make(HTTPRequestParams)
	params["client_id"] = client.Config.ClientID
	params["scope"] = strings.Join(client.Config.Scopes, " ")
	params["response_type"] = "code"
	params["redirect_uri"] = client.Config.RedirectURL
	uri := client.buildURI("https://login.microsoftonline.com/common/oauth2/v2.0/authorize", params)
	return uri
}

func (client *Client) redeemCodeForAccessToken(code string) (*LoginRedeemCodeResponse, error) {
	params := make(HTTPRequestParams)
	params["client_id"] = client.Config.ClientID
	params["redirect_uri"] = client.Config.RedirectURL
	params["client_secret"] = client.Config.ClientSecret
	params["code"] = code
	params["grant_type"] = "authorization_code"
	status, resp, err := client.httpPostForm("https://login.microsoftonline.com/common/oauth2/v2.0/token", params)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		if status == http.StatusUnauthorized {
			return nil, errors.New("verify you're using the client secret's value (not ID) and the API permissions are set correctly")
		}
		return nil, client.handleResponseError(status, resp)
	}
	var json LoginRedeemCodeResponse
	if err := UnmarshalJSON(&json, resp); err != nil {
		return nil, err
	}
	client.UpdateSecretStore(&json)
	return &json, nil
}

func (client *Client) ShouldRenewAccessToken() bool {
	now := time.Now()
	diff := now.Sub(client.Config.Expiry)
	return diff.Minutes() > -30
}

func (client *Client) RenewAccessToken() (*LoginRedeemCodeResponse, error) {
	params := make(HTTPRequestParams)
	params["client_id"] = client.Config.ClientID
	params["redirect_uri"] = client.Config.RedirectURL
	params["client_secret"] = client.Config.ClientSecret
	params["refresh_token"] = client.Config.RefreshToken
	params["grant_type"] = "refresh_token"
	status, resp, err := client.httpPostForm("https://login.microsoftonline.com/common/oauth2/v2.0/token", params)
	if err != nil {
		return nil, err
	}
	if status != http.StatusOK {
		return nil, client.handleResponseError(status, resp)
	}
	var json LoginRedeemCodeResponse
	if err := UnmarshalJSON(&json, resp); err != nil {
		return nil, err
	}
	client.UpdateSecretStore(&json)
	return &json, nil
}

func (client *Client) expectCode() string {
	httpServer := &http.Server{
		Addr:         "0.0.0.0:53682",
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
	}
	ctx, cancel := context.WithCancel(context.Background())

	code := ""
	var handleCall = func(w http.ResponseWriter, r *http.Request) {
		s := r.URL.Query().Get("code")
		if s != "" {
			html := loginHTMLResponseHeader + loginHTMLResponseOK + loginHTMLResponseFooter
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusOK)
			w.Write([]byte(html))
			if f, ok := w.(http.Flusher); ok {
				f.Flush()
			}
			code = s
			cancel()
		} else {
			html := loginHTMLResponseHeader + loginHTMLResponseNotFound + loginHTMLResponseFooter
			w.Header().Set("Content-Type", "text/html")
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte(html))
		}
	}
	http.HandleFunc("/", handleCall)

	go func() error {
		return httpServer.ListenAndServe()
	}()
	<-ctx.Done()
	httpServer.Shutdown(context.Background())

	return code
}

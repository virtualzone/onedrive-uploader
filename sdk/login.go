package sdk

import (
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
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

type SecretStore struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	Expiry       time.Time `json:"expiry"`
}

func (client *Client) Login() error {
	code := client.expectCode()
	grant, _ := client.redeemCodeForAccessToken(code)
	if grant.AccessToken == "" {
		return errors.New("received empty access token")
	}
	return nil
}

func (client *Client) UpdateSecretStore(grant *LoginRedeemCodeResponse) error {
	expiry := time.Now().Add(time.Second * time.Duration(grant.ExpiresIn))
	store := &SecretStore{
		AccessToken:  grant.AccessToken,
		RefreshToken: grant.RefreshToken,
		Expiry:       expiry,
	}
	data, err := json.Marshal(store)
	if err != nil {
		return err
	}
	if err := ioutil.WriteFile(client.Config.SecretStore, data, 0600); err != nil {
		return err
	}
	client.SecretStore = store
	return nil
}

func (client *Client) ReadSecretStore() error {
	data, err := ioutil.ReadFile(client.Config.SecretStore)
	if err != nil {
		return err
	}
	var store SecretStore
	if err := UnmarshalJSON(&store, data); err != nil {
		return err
	}
	client.SecretStore = &store
	return nil
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
	_, resp, err := client.httpPostForm("https://login.microsoftonline.com/common/oauth2/v2.0/token", params)
	if err != nil {
		return nil, err
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
	diff := now.Sub(client.SecretStore.Expiry)
	return diff.Minutes() > -30
}

func (client *Client) RenewAccessToken() (*LoginRedeemCodeResponse, error) {
	params := make(HTTPRequestParams)
	params["client_id"] = client.Config.ClientID
	params["redirect_uri"] = client.Config.RedirectURL
	params["client_secret"] = client.Config.ClientSecret
	params["refresh_token"] = client.SecretStore.RefreshToken
	params["grant_type"] = "refresh_token"
	_, resp, err := client.httpPostForm("https://login.microsoftonline.com/common/oauth2/v2.0/token", params)
	if err != nil {
		return nil, err
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
	code := ""
	var handleCall = func(w http.ResponseWriter, r *http.Request) {
		s := r.URL.Query().Get("code")
		if s != "" {
			w.WriteHeader(http.StatusOK)
			code = s
			httpServer.Shutdown(context.TODO())
		} else {
			w.WriteHeader(http.StatusNotFound)
		}
	}
	http.HandleFunc("/", handleCall)
	httpServer.ListenAndServe()

	return code
}

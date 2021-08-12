package client

import (
	"encoding/json"
	"fmt"
	"github.com/google/go-querystring/query"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

var PasswordEncryptionTypes = []string{"PLAIN", "SHA", "LDAP_SHA", "MD4", "MD5"}

type Client struct {
	HostUrl    string
	HTTPClient HttpClient
	Token      string
}

type HttpClient interface {
	Do(req *http.Request) (*http.Response, error)
}

type ReportPortalClientConfig struct {
	Username, Password, Host string
}

type uiAuthResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    uint   `json:"expires_in"`
	Scope        string `json:"ui"`
	Jti          string `json:"jti"`
}

func NewClient(config *ReportPortalClientConfig, httpClient HttpClient) (*Client, error) {
	if len(config.Host) <= 0 || len(config.Username) <= 0 || len(config.Password) <= 0 {
		return nil, fmt.Errorf("host, username and password are required parameters")
	}

	if httpClient == nil {
		log.Print("ReportPortal client using the default HTTPClient: Timeout 10s")
		httpClient = &http.Client{Timeout: 10 * time.Second}
	}

	c := Client{
		HTTPClient: httpClient,
		HostUrl:    config.Host,
	}

	uiAccessToken, err := c.getUiAccessToken(&config.Username, &config.Password)
	if err != nil {
		return nil, err
	}
	c.Token = *uiAccessToken

	return &c, err
}

func (c *Client) getUiAccessToken(username, password *string) (*string, error) {
	data := url.Values{}
	data.Set("grant_type", "password")
	data.Set("username", *username)
	data.Set("password", *password)

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/uat/sso/oauth/token", c.HostUrl), strings.NewReader(data.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", "Basic dWk6dWltYW4=")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	body, err := c.doRequest(req)
	if err != nil {
		return nil, err
	}

	ar := uiAuthResponse{}
	err = json.Unmarshal(body, &ar)
	if err != nil {
		return nil, err
	}

	return &ar.AccessToken, nil
}

func (c *Client) doRequest(req *http.Request) ([]byte, error) {
	if req.Header.Get("Authorization") == "" {
		req.Header.Add("Authorization", fmt.Sprintf("Bearer %s", c.Token))
	}

	res, err := c.HTTPClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	if res.StatusCode < http.StatusOK || res.StatusCode > http.StatusAlreadyReported {
		return nil, fmt.Errorf("status: %d, body: %s", res.StatusCode, body)
	}

	return body, err
}

func parseQuery(i interface{}) string {
	if i == nil {
		return ""
	}

	v, _ := query.Values(i)
	return v.Encode()
}

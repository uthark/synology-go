package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

func New(endpoint string) *SynologyClient {
	client := http.Client{}

	return &SynologyClient{
		client:   client,
		Endpoint: endpoint,
	}

}

type SynologyClient struct {
	Endpoint  string
	client    http.Client
	sessionID *string
}

func (c *SynologyClient) Logout() error {
	endpoint := "%s/webapi/auth.cgi?api=SYNO.API.Auth&method=logout&version=1&_sid=%s&session=dsm_info"
	_, err := c.client.Get(fmt.Sprintf(endpoint, c.Endpoint, url.QueryEscape(*c.sessionID)))
	return err
}

func (c *SynologyClient) Auth(login, password string) (*AuthResponse, error) {
	endpoint := "%s/webapi/auth.cgi?api=SYNO.API.Auth&method=login&version=3&account=%s&passwd=%s&format=sid&session=dsm_info"
	resp, err := c.client.Get(fmt.Sprintf(endpoint, c.Endpoint, url.QueryEscape(login), url.QueryEscape(password)))
	if err != nil {
		return nil, err
	}
	fmt.Println(resp.StatusCode, resp.Status)
	decoder := json.NewDecoder(resp.Body)
	result := &AuthResponse{}

	err = decoder.Decode(result)
	if err != nil {
		return nil, err
	}

	if err == nil {
		fmt.Println(result.AuthData.SessionID)
		c.sessionID = &result.AuthData.SessionID
	}
	fmt.Println((*c).sessionID)
	return result, err
}

type AuthData struct {
	SessionID    string `json:"sid"`
	DeviceID     string `json:"did"`
	IsPortalPort bool   `json:"is_portal_port"`
}

type AuthResponse struct {
	AuthData AuthData `json:"data"`
	Success  bool     `json:"success"`
}

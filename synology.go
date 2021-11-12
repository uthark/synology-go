package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	"github.com/joho/godotenv"
)

func main() {

	err := godotenv.Load()
	if err != nil {
		panic(err)
	}
	endpoint := os.Getenv("SYNOLOGY_HOST")
	client := New(endpoint)
	response, err := client.Auth(os.Getenv("SYNOLOGY_LOGIN"), os.Getenv("SYNOLOGY_PASSWORD"))
	if err != nil {
		fmt.Println(err)
		return
	}

	client.Info()

	hostAPI := VirtualizationAPIHost{
		Endpoint: endpoint,
		Session:  response.AuthData.SessionID,
	}
	list, err := hostAPI.List()
	fmt.Println(list, err)
}

type SynologyClient struct {
	Endpoint  string
	client    http.Client
	sessionID string
}

type VirtualizationAPIHost struct {
	Session  string
	Endpoint string
}

type Host struct {
	ID             string `json:"host_id"`
	Name           string `json:"host_name"`
	Status         string `json:"status"`
	TotalCPUCore   int    `json:"total_cpu_core"`
	TotalRAMSizeMB int    `json:"total_ram_size"`
	FreeCPUCore    int    `json:"free_cpu_core"`
	FreeRAMSizeMB  int    `json:"free_ram_size"`
}

type ListData struct {
	Hosts []Host `json:"hosts"`
}

type ListResponse struct {
	ListData ListData `json:"data"`
	Success  bool     `json:"success"`
}

func (h VirtualizationAPIHost) List() (*ListResponse, error) {
	client := http.Client{}

	endpoint := "%s/webapi/entry.cgi?api=SYNO.Virtualization.API.Host&method=list&version=1&_sid=%s"
	resp, err := client.Get(fmt.Sprintf(endpoint, h.Endpoint, h.Session))
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(resp.Body)
	result := &ListResponse{}
	err = decoder.Decode(result)

	return result, err
}

func (c SynologyClient) Auth(login, password string) (*AuthResponse, error) {
	endpoint := "%s/webapi/auth.cgi?api=SYNO.API.Auth&method=login&version=3&account=%s&passwd=%s&format=sid&session=dsm_info"
	resp, err := c.client.Get(fmt.Sprintf(endpoint, c.Endpoint, url.QueryEscape(login), url.QueryEscape(password)))
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(resp.Body)
	result := &AuthResponse{}

	err = decoder.Decode(result)

	if err != nil {
		c.sessionID = result.AuthData.SessionID
	}
	return result, err
}

type APIInfo struct {
	MinVersion int `json:"minVersion"`
	MaxVersion int `json:"maxVersion"`
	RequestFormat string `json:"requestFormat"`
	Path string `json:"path"`
}

type InfoData map[string]APIInfo

type InfoResponse struct {
	Data InfoData `json:"data"`
	Success bool `json:"success"`
}

func (c SynologyClient) Info() (*InfoResponse, error) {
	endpoint := "%s/webapi/entry.cgi?api=SYNO.API.Info&version=1&method=query"
	resp, err := c.client.Get(fmt.Sprintf(endpoint, c.Endpoint))
	if err != nil {
		return nil, err
	}
	bytes, err := ioutil.ReadAll(resp.Body)
	fmt.Println(string(bytes))
	result := &InfoResponse{}
	err = json.Unmarshal(bytes, resp)
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

func New(endpoint string) *SynologyClient {
	client := http.Client{}

	return &SynologyClient{
		client:   client,
		Endpoint: endpoint,
	}

}

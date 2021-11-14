package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
)

type APIInfo struct {
	MinVersion    int    `json:"minVersion"`
	MaxVersion    int    `json:"maxVersion"`
	RequestFormat string `json:"requestFormat"`
	Path          string `json:"path"`
}

type InfoData map[string]APIInfo

type InfoResponse struct {
	Data    InfoData `json:"data"`
	Success bool     `json:"success"`
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

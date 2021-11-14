package main

import (
	bytes2 "bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type VirtualizationAPI struct {
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

type ListHostsData struct {
	Hosts []Host `json:"hosts"`
}

type ListHostsResponse struct {
	Data    ListHostsData `json:"data"`
	Success bool          `json:"success"`
}

func (h VirtualizationAPI) ListHosts() (*ListHostsResponse, error) {
	client := http.Client{}

	endpoint := "%s/webapi/entry.cgi?api=SYNO.Virtualization.API.Host&method=list&version=1&_sid=%s"
	resp, err := client.Get(fmt.Sprintf(endpoint, h.Endpoint, h.Session))
	if err != nil {
		return nil, err
	}
	decoder := json.NewDecoder(resp.Body)
	result := &ListHostsResponse{}
	err = decoder.Decode(result)

	return result, err
}

type DiskController int

const (
	DiskControllerVirtIO DiskController = 1
	DiskControllerIDE    DiskController = 2
	DiskControllerSATA   DiskController = 3
)

type Disk struct {
	ID         string         `json:"vdisk_id"`
	Controller DiskController `json:"controller"`
	Unmap      bool           `json:"unmap"`
	SizeInMB   bool           `json:"vdisk_size"`
}

type NicModel int

const (
	NicModelVirtIO  NicModel = 1
	NicModelE1000   NicModel = 2
	NicModelRTL8139 NicModel = 3
)

type NIC struct {
	ID          string   `json:"vnic_id"`
	Model       NicModel `json:"model"`
	MacAddress  string   `json:"mac"`
	NetworkID   string   `json:"network_id"`
	NetworkName string   `json:"network_name"`
}

type GuestAutorun int

const (
	AutorunOff       GuestAutorun = 0
	AutorunLastState GuestAutorun = 1
	AutorunOn        GuestAutorun = 2
)

type Guest struct {
	Autorun     GuestAutorun `json:"autorun"`
	Description string       `json:"description"`
	ID          string       `json:"guest_id"`
	Name        string       `json:"guest_name"`
	Status      string       `json:"status"`
	StorageID   string       `json:"storage_id"`
	StorageName string       `json:"storage_name"`
	CPUs        int          `json:"vcpu_num"`
	RamInMB     int          `json:"vram_size"`
	Disks       []Disk       `json:"vdisks"`
	NICs        []NIC        `json:"vnics"`
}

type ListGuestsData struct {
	Guests []Guest `json:"guests"`
}

type ListGuestsResponse struct {
	Data    ListGuestsData `json:"data"`
	Success bool           `json:"success"`
}

func (h VirtualizationAPI) ListGuests() (*ListGuestsResponse, error) {
	client := http.Client{}

	endpoint := "%s/webapi/entry.cgi?api=SYNO.Virtualization.API.Guest&method=list&version=1&_sid=%s"
	resp, err := client.Get(fmt.Sprintf(endpoint, h.Endpoint, h.Session))
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(b))
	decoder := json.NewDecoder(bytes2.NewReader(b))
	result := &ListGuestsResponse{}
	err = decoder.Decode(result)

	return result, err
}

type GetGuestRequest struct {
	ID *string `json:"guest_id"`
	Name *string `json:"guest_name"`
	Additional bool  `json:"additional"`
}

type GetGuestResponse struct {
	Data    Guest `json:"data"`
	Success bool           `json:"success"`
}

func (h VirtualizationAPI) GetGuest(req GetGuestRequest) (*GetGuestResponse, error) {
	client := http.Client{}

	var key string
	var val string
	if req.ID != nil {
		key = "guest_id"
		val = *(req.ID)
	} else {
		key = "guest_name"
		val = *(req.Name)
	}
	endpoint := "%s/webapi/entry.cgi?api=SYNO.Virtualization.API.Guest&method=get&version=1&_sid=%s&%s=%s&additional=%t"
	url := fmt.Sprintf(endpoint, h.Endpoint, h.Session, key, val, req.Additional)
	fmt.
		Println(url)
	resp, err := client.Get(url)
	if err != nil {
		return nil, err
	}
	b, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	fmt.Println(string(b))
	decoder := json.NewDecoder(bytes2.NewReader(b))
	result := &GetGuestResponse{}
	err = decoder.Decode(result)

	return result, err
}

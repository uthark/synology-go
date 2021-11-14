package main

import (
	"fmt"
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

	//client.Info()

	v := VirtualizationAPI{
		Endpoint: endpoint,
		Session:  response.AuthData.SessionID,
	}

	guestsResponse, err := v.ListGuests()
	guests := guestsResponse.Data.Guests
	for i := range guests {
		fmt.Println(guests[i].Name, guests[i].Autorun)
	}

	s := "vm03"
	guest, err := v.GetGuest(GetGuestRequest{Name: &s})
	fmt.Println(guest.Data.Status)

	client.Logout()
}

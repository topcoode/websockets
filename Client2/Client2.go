package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type User struct {
	Name     string
	Username string
	Status   string
	Message  string
	Online   int
}

func main() {
	// Send a GET request to the server
	// resp, err := http.Get("http://localhost:8080/client2")
	// if err != nil {
	// 	fmt.Println("Error sending GET request:", err)
	// 	return
	// }
	req, err := http.NewRequest("GET", "http://localhost:8080/client1", nil) //post
	if err != nil {
		fmt.Println("error in link", err)
	}

	// Read the response body
	defer req.Body.Close()
	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		fmt.Println("Error reading response body:", err)
		return
	}

	// Print the response body
	fmt.Println("Response:", string(body))
	var userdata User
	r, err := json.Marshal(userdata)
	if err != nil {
		fmt.Println("error on marshalling..........", err)

	}
	fmt.Println(string(r))
}

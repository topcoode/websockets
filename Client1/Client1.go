package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

var clients = make(map[*websocket.Conn]bool) // Connected clients
var upgrader = websocket.Upgrader{}

type User struct {
	Name     string
	Username string
	Status   string
	Message  string
	Online   string
}

func main() {
	user := User{}
	fmt.Println("client One is active..")
	fmt.Print("Enter your name : ")
	fmt.Scanln(&user.Name)
	fmt.Print("Enter your username : ")
	fmt.Scanln(&user.Username)
	fmt.Print("Enter your status : ")
	fmt.Scanln(&user.Status)
	fmt.Print("Enter your message : ")
	fmt.Scanln(&user.Message)
	fmt.Print("Enter your online : ")
	fmt.Scanln(&user.Online)

	UserJson, err := json.Marshal(user)
	if err != nil {
		fmt.Println("error on marshalling..........", err)

	}

	req, err := http.NewRequest("POST", "http://localhost:8080/client1", bytes.NewBuffer(UserJson)) //post
	fmt.Println(req)
	if err != nil {
		fmt.Println("error in link", err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()

}

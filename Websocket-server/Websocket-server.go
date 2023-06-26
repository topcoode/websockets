package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

type Client struct {
	conn *websocket.Conn
	send chan []byte
}

var (
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			// Allow all connections
			return true
		},
	}

	clients    = make(map[*Client]bool)
	broadcast  = make(chan []byte)
	register   = make(chan *Client)
	unregister = make(chan *Client)
)

func main() {
	p := 8081
	fmt.Println("the value of p------------------>", p)
	fmt.Printf("Starting websocket echo server on port %d", p)

	http.HandleFunc("/", echo)
	if err := http.ListenAndServe(fmt.Sprintf(":%d", p), nil); err != nil {
		fmt.Printf("Error while starting to listen: %#v", err)
	}
}

func echo(w http.ResponseWriter, r *http.Request) {
	// Start handshaking
	c, err := upgrader.Upgrade(w, r, nil)
	fmt.Println("the value c ------------------>", c)
	if err != nil {
		fmt.Printf("Upgrading error: %#v\n", err)
		return
	}
	defer c.Close()
	// End

	fmt.Println("Success to handshake with client")

	// Start bidirectional messages
	for {
		mt, message, err := c.ReadMessage()
		fmt.Println("the value of mt --------------------->", mt)
		if err != nil {
			fmt.Printf("Reading error: %#v\n", err)
			break
		}
		fmt.Printf("recv: message %q", message)
		if err := c.WriteMessage(mt, message); err != nil {
			fmt.Printf("Writing error: %#v\n", err)
			break
		}
	}
	// End bidirectional messages
}
func runWebSocketHub() {
	for {
		select {
		case client := <-register:
			clients[client] = true
		case client := <-unregister:
			if _, ok := clients[client]; ok {
				delete(clients, client)
				close(client.send)
			}
		case message := <-broadcast:
			for client := range clients {
				select {
				case client.send <- message:
				default:
					close(client.send)
					delete(clients, client)
				}
			}
		}
	}
}

func (c *Client) readPump() {
	defer func() {
		unregister <- c
		c.conn.Close()
	}()

	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}

		broadcast <- message
	}
}

func (c *Client) writePump() {
	defer func() {
		c.conn.Close()
	}()

	for message := range c.send {
		err := c.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			break
		}
	}
}

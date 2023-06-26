package main

import (
	"context"
	"fmt"

	"github.com/gorilla/websocket"
)

func main() {
	ctx := context.Background()
	endpointUrl := "ws://localhost:8081"
	dialer := websocket.Dialer{}
	c, _, err := dialer.DialContext(ctx, endpointUrl, nil)
	fmt.Println("ctx--------------->", ctx)
	fmt.Println("endpoint--------------->", endpointUrl)
	fmt.Println("c----------------->", c)

	if err != nil {
		fmt.Println(err)
	}
	defer c.Close()
	done := make(chan struct{})
	go func() {
		defer c.Close()
		defer close(done)

		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				fmt.Println(err)
				return
			}
			fmt.Printf("recv: %s\n", message)
		}
	}()

}

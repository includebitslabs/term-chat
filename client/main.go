// socket-client project main.go
package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"strings"

	"golang.org/x/net/websocket"
)

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = "9988"
	SERVER_TYPE = "tcp"
)

func main() {
	reader := bufio.NewReader(os.Stdin)
	fmt.Printf("enter client name : ")
	clientName, _ := reader.ReadString('\n')
	clientName = strings.TrimSpace(clientName)
	origin := "http://localhost/" + clientName
	url := "ws://localhost:9988"
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}
	defer ws.Close()

	go func() {
		for {
			var msg = make([]byte, 512)
			n, err := ws.Read(msg)
			if err != nil {
				log.Println("Error reading from server:", err)
				return
			}
			fmt.Printf("\nReceived: %s\n", msg[:n])
		}
	}()

	// Main loop to handle user input
	for {
		fmt.Print("\nText to send: ")
		text, _ := reader.ReadString('\n')
		text = strings.TrimSpace(text) // Remove the newline character

		// Send the message to the WebSocket server
		_, err := ws.Write([]byte(text))
		if err != nil {
			log.Fatal(err)
		}
	}
}

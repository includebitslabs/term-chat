// socket-server project main.go
package main

import (
	"fmt"
	"io"
	"net/http"

	"golang.org/x/net/websocket"
)

const (
	SERVER_HOST = "localhost"
	SERVER_PORT = ":9988"
	SERVER_TYPE = "tcp"
)

type Server struct {
	connections map[*websocket.Conn]bool
}

func NewServer() *Server {
	return &Server{
		connections: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) handleWS(ws *websocket.Conn) {
	fmt.Println("\nnew incoming connection from client :", ws.RemoteAddr())
	s.connections[ws] = true

	s.readLoop(ws)
}

func (s *Server) readLoop(ws *websocket.Conn) {
	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				break
			}
			fmt.Println("read error", err)
			continue
		}
		msg := buf[:n]
		s.broadcast(ws, msg)
	}
}

func (s *Server) broadcast(origin *websocket.Conn, b []byte) {
	senderAddr := origin.RemoteAddr().String()
	fullMsg := fmt.Sprintf("%s: %s", senderAddr, string(b))

	for ws := range s.connections {
		if ws == origin {
			continue // Skip sending the message to the original sender
		}
		go func(ws *websocket.Conn) {
			if _, err := ws.Write([]byte(fullMsg)); err != nil {
				fmt.Println("write error :: ", err)
			}
		}(ws)
	}
}

func main() {
	server := NewServer()
	fmt.Printf("server starting at port %s", SERVER_PORT)
	http.Handle("/", websocket.Handler(server.handleWS))
	http.ListenAndServe(SERVER_PORT, nil)
}

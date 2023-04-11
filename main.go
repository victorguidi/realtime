package main

// TODO: add a mutex to protect the map
// TODO: find a way to add TLS support
// TODO: add a way to close the connection from the server side
// TODO: add a way to authenticate the client

import (
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

type Server struct {
	// serverId string
	// serverToken string
	// private bool
	conns map[*websocket.Conn]bool
}

func NewServer() *Server {
	return &Server{
		conns: make(map[*websocket.Conn]bool),
	}
}

func (s *Server) handleWSOrderbook(ws *websocket.Conn) {
	fmt.Println("new incoming connection from client:", ws.RemoteAddr())

	for {
		payload := fmt.Sprintf("orderbook data -> %d\n", time.Now().UnixNano())
		ws.Write([]byte(payload))
		time.Sleep(2 * time.Second)
	}
}

func (s *Server) handleWSOrderbookWithAuth(ws *websocket.Conn) {
	fmt.Println("ws.Config().Header:", ws.Config().Header)
	userToken := ws.Config().Header.Get("User-Token")
	if !s.validateUser(userToken) {
		fmt.Println("invalid user token", userToken)
		ws.Close()
		return
	}
	fmt.Println("new incoming connection from client:", ws.RemoteAddr())
	for {
		payload := fmt.Sprintf("orderbook data -> %d\n", time.Now().UnixNano())
		ws.Write([]byte(payload))
		time.Sleep(2 * time.Second)
	}
}

func (s *Server) validateUser(userToken string) bool {
	if userToken == "123" {
		return true
	}
	return false
}

func (s *Server) handleWS(ws *websocket.Conn) {
	fmt.Println("new incoming connection from client:", ws.RemoteAddr())

	// We should use a mutex here to protect the map
	s.conns[ws] = true

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
			fmt.Println("error reading from client:", err)
			continue
		}
		msg := buf[:n]

		s.broadcast(msg)
	}
}

func (s *Server) broadcast(b []byte) {
	for ws := range s.conns {
		go func(ws *websocket.Conn) {
			if _, err := ws.Write(b); err != nil {
				fmt.Println("error writing to client:", err)
			}
		}(ws)
	}
}

func main() {
	server := NewServer()
	http.Handle("/ws", websocket.Handler(server.handleWS))
	http.Handle("/ws/orderbook", websocket.Handler(server.handleWSOrderbook))
	http.Handle("/ws/auth/orderbook", websocket.Handler(server.handleWSOrderbookWithAuth))
	http.ListenAndServe(":8080", nil)

}

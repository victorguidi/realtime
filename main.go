package main

// TODO: add a mutex to protect the map
// TODO: find a way to add TLS support
// TODO: implement sessions

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"

	"golang.org/x/net/websocket"
)

type Server struct {
	// serverId string
	serverToken string
	// private bool
	conns map[*websocket.Conn]bool
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func NewServer() *Server {
	return &Server{
		conns:       make(map[*websocket.Conn]bool),
		serverToken: "hello",
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
	h := ws.Request().Header.Get("Sec-Websocket-Protocol")
	if h == s.serverToken {
		buf := make([]byte, 1024)
		ws.Write([]byte("Connection Open, validating user"))
		n, err := ws.Read(buf)
		if err != nil {
			ws.Write([]byte("The json was not correct"))
			ws.Close()
			return
		}

		user := User{}
		err = json.Unmarshal(buf[:n], &user)
		if err != nil {
			ws.Write([]byte("The json was not correct"))
			ws.Close()
			return
		}
		if user.Password == "123" {
			for {
				n, err := ws.Read(buf)
				if err != nil {
					ws.Write([]byte("Something went wrong with the message"))
					ws.Close()
					return
				}
				if string(buf[:n]) == "1" {
					ws.Write([]byte("Thanks for connecting"))
					ws.Close()
					return
				}
				payload := fmt.Sprintf("orderbook data -> %d\n", time.Now().UnixNano())
				ws.Write([]byte(payload))
				time.Sleep(2 * time.Second)
			}
		} else {
			ws.Write([]byte("404"))
		}
	} else {
		ws.Write([]byte("404"))
	}
	ws.Close()
	return
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

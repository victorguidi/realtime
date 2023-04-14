package main

// TODO: add a mutex to protect the map
// TODO: find a way to add TLS support (Done?)
// TODO: implement sessions
// TODO: Implement the api server, in order to login the users and save the sessions

import (
	// "crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	// "os"
	// "time"

	"golang.org/x/net/websocket"
)

type Server struct {
	// serverId string
	serverToken string
	// private bool
	conns map[User]bool
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	connInfo *websocket.Conn
}

type Message struct {
	to  string
	msg []byte
}

func NewServer() *Server {
	return &Server{
		conns:       make(map[User]bool),
		serverToken: "hello",
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

		user := User{
			connInfo: ws,
		}
		err = json.Unmarshal(buf[:n], &user)
		if err != nil {
			ws.Write([]byte("The json was not correct"))
			ws.Close()
			return
		}
		if user.Password == "123" {
			// We should use a mutex here to protect the map
			s.conns[user] = true

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

				m := Message{
					to:  "test",
					msg: buf[:n],
				}

				s.readLoop(&m)
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

func (s *Server) readLoop(msg *Message) {
	for c := range s.conns {
		if c.Username == msg.to {
			c.connInfo.Write(msg.msg)
		}
	}
}

func (s *Server) broadcast(b []byte) {
	for ws := range s.conns {
		go func(ws *websocket.Conn) {
			if _, err := ws.Write(b); err != nil {
				fmt.Println("error writing to client:", err)
			}
		}(ws.connInfo)
	}
}

func main() {

	server := NewServer()
	// http.Handle("/ws/orderbook", websocket.Handler(server.handleWSOrderbook))
	http.Handle("/wss/auth/orderbook", websocket.Handler(server.handleWSOrderbookWithAuth))
	log.Fatal(http.ListenAndServeTLS(":8080", "./selfCertificate/server.crt", "./selfCertificate/server.key", nil))

}

// func (s *Server) handleWSOrderbook(ws *websocket.Conn) {
// 	fmt.Println("new incoming connection from client:", ws.RemoteAddr())
//
// 	for {
// 		payload := fmt.Sprintf("orderbook data -> %d\n", time.Now().UnixNano())
// 		ws.Write([]byte(payload))
// 		time.Sleep(2 * time.Second)
// 	}
// }

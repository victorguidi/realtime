package main

import (
	// "crypto/x509"
	"encoding/json"
	"fmt"
	"log"
	"time"

	// "fmt"
	"net/http"

	// "os"
	// "time"

	"golang.org/x/net/websocket"
)

type Server struct {
	listenAddr string
	store      Storage
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	connInfo *websocket.Conn
	Sessions []int `json:"sessions"`
}

type Message struct {
	to  string
	msg []byte
}

type Session struct {
	ID           int    `json:"id"`
	SessionToken string `json:"sessionToken"`
	conns        map[string]bool
	state        bool
	Created_at   string `json:"created_at"`
	Updated_at   string `json:"updated_at"`
	// private     bool
	SessionName string `json:"sessionName"`
}

func NewServer(listenAddr string, store Storage) *Server {
	return &Server{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *Server) handleGetOpenSessions(w http.ResponseWriter, r *http.Request) {
	type Response struct {
		Sessions []*Session `json:"sessions"`
	}
	var response Response
	sessions, err := s.store.GetAllSessions()
	if err != nil {
		log.Fatal(err)
	}
	response.Sessions = sessions
	json.NewEncoder(w).Encode(&response)
}

func (s *Server) handleRegisterNewUser(w http.ResponseWriter, r *http.Request) {

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	err = s.store.CreateUser(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	json.NewEncoder(w).Encode(user)
}

func (s *Server) handleCreateNewSession(w http.ResponseWriter, r *http.Request) {

	type Response struct {
		SessionId    string `json:"sessionId"`
		SessionToken string `json:"sessionToken"`
	}

	var response Response
	err := json.NewDecoder(r.Body).Decode(&response)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	session := Session{
		SessionToken: response.SessionToken,
		conns:        make(map[string]bool),
		SessionName:  response.SessionId,
		state:        true,
	}

	err = s.store.CreateSession(&session)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	json.NewEncoder(w).Encode(&session)
}

func (s *Server) handleUserSessions(w http.ResponseWriter, r *http.Request) {

	type Response struct {
		SessionId int `json:"sessionId"`
		UserId    int `json:"userId"`
	}

	var resp Response
	err := json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
	}

	err = s.store.AddSessionToUser(resp.UserId, resp.SessionId)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
	}

	json.NewEncoder(w).Encode(&resp)
}

func (s *Server) handleGetUsers(w http.ResponseWriter, r *http.Request) {

	users, err := s.store.GetAllUsers()
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusBadGateway)
	}
	json.NewEncoder(w).Encode(&users)

}

func (s *Server) handleUpgradeToWsSession(ws *websocket.Conn) {
	// TODO: Spin up a new go connection that will handle new connections
	fmt.Println("new incoming connection from client:", ws.RemoteAddr())

	// if ws.Request().Header.Get("sessionId") != "1" {
	// 	log.Fatal("What")
	// }
	//
	// if ws.Request().Header.Get("sessionToken") != "123" {
	// 	log.Fatal("What")
	// }
	//
	// if ws.Request().Header.Get("Upgrade") != "websocket" {
	// 	log.Fatal("What")
	// }

	for {
		payload := fmt.Sprintf("orderbook data -> %d\n", time.Now().UnixNano())
		ws.Write([]byte(payload))
		time.Sleep(2 * time.Second)
	}

	// type Response struct {
	// 	SessionId    int    `json:"sessionId"`
	// 	SessionToken string `json:"sessionToken"`
	// }
	//
	// var resp Response
	// err := json.NewDecoder(r.Body).Decode(&resp)
	// if err != nil {
	// 	w.WriteHeader(http.StatusBadGateway)
	// }
	//
	// session, err := s.store.GetSession(resp.SessionId)
}

func (s *Server) handleWSConnection(ws *websocket.Conn) {

}

// func (s *Session) broadcast(b []byte) {
// 	for ws := range s.conns {
// 		go func(ws *websocket.Conn) {
// 			if _, err := ws.Write(b); err != nil {
// 				fmt.Println("error writing to client:", err)
// 			}
// 		}(ws.connInfo)
// 	}
// }

// func (s *Server) handleWSOrderbook(ws *websocket.Conn) {
// 	fmt.Println("new incoming connection from client:", ws.RemoteAddr())
//
// 	for {
// 		payload := fmt.Sprintf("orderbook data -> %d\n", time.Now().UnixNano())
// 		ws.Write([]byte(payload))
// 		time.Sleep(2 * time.Second)
// 	}
// }

// func (s *Server) handleWSOrderbookWithAuth(ws *websocket.Conn) {
// 	h := ws.Request().Header.Get("Sec-Websocket-Protocol")
// 	if h == s.serverToken {
// 		buf := make([]byte, 1024)
// 		ws.Write([]byte("Connection Open, validating user"))
// 		n, err := ws.Read(buf)
// 		if err != nil {
// 			ws.Write([]byte("The json was not correct"))
// 			ws.Close()
// 			return
// 		}
//
// 		user := User{
// 			connInfo: ws,
// 		}
// 		err = json.Unmarshal(buf[:n], &user)
// 		if err != nil {
// 			ws.Write([]byte("The json was not correct"))
// 			ws.Close()
// 			return
// 		}
// 		if user.Password == "123" {
// 			// We should use a mutex here to protect the map
// 			s.conns[user] = true
//
// 			for {
//
// 				n, err := ws.Read(buf)
// 				if err != nil {
// 					ws.Write([]byte("Something went wrong with the message"))
// 					ws.Close()
// 					return
// 				}
// 				if string(buf[:n]) == "1" {
// 					ws.Write([]byte("Thanks for connecting"))
// 					ws.Close()
// 					return
// 				}
//
// 				m := Message{
// 					to:  "test",
// 					msg: buf[:n],
// 				}
//
// 				s.readLoop(&m)
// 			}
// 		} else {
// 			ws.Write([]byte("404"))
// 		}
// 	} else {
// 		ws.Write([]byte("404"))
// 	}
// 	ws.Close()
// 	return
// }

// func (s *Server) readLoop(msg *Message) {
// 	for c := range s.conns {
// 		if c.Username == msg.to {
// 			c.connInfo.Write(msg.msg)
// 		}
// 	}
// }

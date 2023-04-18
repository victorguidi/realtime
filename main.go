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
	conns    map[User]bool
	Sessions []Session `json:"Sessions"`
	Users    []User
}

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type WUser struct {
	Username string `json:"username"`
	Password string `json:"password"`
	connInfo *websocket.Conn
}

type Message struct {
	to  string
	msg []byte
}

type Session struct {
	SessionToken string `json:"sessionToken"`
	conns        map[WUser]bool
	state        bool
	// private     bool
	SessionId string `json:"sessionId"`
}

func NewServer() *Server {
	return &Server{
		conns:       make(map[User]bool),
		serverToken: "hello",
		Sessions:    make([]Session, 0),
	}
}

func (s *Server) handleGetOpenSessions(w http.ResponseWriter, r *http.Request) {
	for _, session := range s.Sessions {
		if session.state == true {
			s.Sessions = append(s.Sessions, session)
		}
	}
	json.NewEncoder(w).Encode(&s.Sessions)
}

func (s *Server) handleRegisterNewUser(w http.ResponseWriter, r *http.Request) {

	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
	}

	s.Users = append(s.Users, user)
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
		conns:        make(map[WUser]bool),
		SessionId:    response.SessionId,
		state:        true,
	}

	s.Sessions = append(s.Sessions, session)

	json.NewEncoder(w).Encode(&session)
}

func (s *Server) handleAddUsersToSession(w http.ResponseWriter, r *http.Request) {

	type Response struct {
		SessionId string `json:"sessionId"`
		Username  string `json:"username"`
		Password  string `json:"password"`
	}

	var resp Response
	err := json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
	}

	user := WUser{
		Username: resp.Username,
		Password: resp.Password,
	}

	for _, session := range s.Sessions {
		if session.SessionId == resp.SessionId {
			session.conns[user] = true
		}
	}

	json.NewEncoder(w).Encode(&user)
}

func (s *Server) handleUpgradeToWsSession(w http.ResponseWriter, r *http.Request) {

	if r.Header.Get("Upgrade") != "websocket" {
		http.Error(w, "Upgrade to Websocket required", http.StatusBadRequest)
	}

	hijack, ok := w.(http.Hijacker)
	if !ok {
		http.Error(w, "Hijacking is not supported", http.StatusInternalServerError)
	}
	conn, _, err := hijack.Hijack()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	type Response struct {
		SessionId string `json:"sessionId"`
		Username  string `json:"username"`
		Password  string `json:"password"`
	}

	var resp Response
	err = json.NewDecoder(r.Body).Decode(&resp)
	if err != nil {
		w.WriteHeader(http.StatusBadGateway)
	}

	wsHandler := func(ws *websocket.Conn) {
		for session := range s.Sessions {
			ws.Write([]byte("Connection established"))
			_ = session
		}
	}
	websocket.Handler(wsHandler).ServeHTTP(w, r)

	conn.Close()
}

func (s *Session) broadcast(b []byte) {
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
	// http.Handle("/wss/auth/orderbook", websocket.Handler(server.handleWSOrderbookWithAuth))
	http.HandleFunc("/api/getSessions", server.handleGetOpenSessions)
	http.HandleFunc("/api/registerUser", server.handleRegisterNewUser)
	http.HandleFunc("/api/registerSession", server.handleCreateNewSession)
	http.HandleFunc("/api/registerUserToSession", server.handleAddUsersToSession)
	http.HandleFunc("/wss/login", server.handleUpgradeToWsSession)
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

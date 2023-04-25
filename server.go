package main

import (
	// "crypto/x509"
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"io"
	"log"
	"net/http"
	"strconv"
)

type Server struct {
	listenAddr string
	store      Storage
	sessions   []*Session
}

type User struct {
	ID       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	connInfo *websocket.Conn
	Sessions []int `json:"sessions"`
}

type SessionUser struct {
	User_Id  int `json:"user_id"`
	connInfo *websocket.Conn
}

type Message struct {
	User_id int    `json:"from"`
	Message string `json:"message"`
}

type Session struct {
	ID           int    `json:"id"`
	SessionToken string `json:"sessionToken"`
	conns        map[*SessionUser]bool
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
		sessions:   make([]*Session, 0),
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
		conns:        make(map[*SessionUser]bool),
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

func (s *Server) handleWs(ws *websocket.Conn) {
	// TODO: Spin up a new go connection that will handle new connections
	fmt.Println("new incoming connection from client:", ws.RemoteAddr())

	h := ws.Request().Header.Get("Sec-Websocket-Protocol")
	id, err := strconv.Atoi(h)
	if err != nil {
		ws.Write([]byte("Something went wrong reading the session Id"))
		ws.Close()
	}

	var session *Session
	if len(s.sessions) == 0 {
		session, err = s.store.GetSession(id)
		if err != nil {
			ws.Write([]byte("There is no session with this ID"))
			ws.Close()
		}
		session.conns = make(map[*SessionUser]bool)
		session.state = true
		s.sessions = append(s.sessions, session)
	} else {
		for _, s := range s.sessions {
			if s.ID == id {
				session = s
			}
		}
	}

	if session.authenticateUser(ws) {
		session.readLoop(ws)
	} else {
		for user := range session.conns {
			if user.connInfo == ws {
				delete(session.conns, user)
			}
		}
		ws.Close()
	}
}

func (s *Session) authenticateUser(ws *websocket.Conn) bool {
	type Response struct {
		SessionToken string `json:"sessionToken"`
		User_Id      int    `json:"id"`
	}
	var r Response
	user := &SessionUser{
		User_Id:  r.User_Id,
		connInfo: ws,
	}
	s.conns[user] = true
	for {
		if err := websocket.JSON.Receive(ws, &r); err != nil {
			if err == io.EOF {
				fmt.Println("Client disconnected")
				delete(s.conns, user)
				return false
			}
			fmt.Println("error reading from client:", err)
			delete(s.conns, user)
			continue
		}
		if r.SessionToken == s.SessionToken {
			ws.Write([]byte("User authenticated to session"))
			break
		}
	}
	return true
}

func (s *Session) readLoop(ws *websocket.Conn) {
	var r Message

	for {
		if err := websocket.JSON.Receive(ws, &r); err != nil {
			if err == io.EOF {
				fmt.Println("Client disconnected")
				ws.Close()
				return
			}
			fmt.Println("error reading from client:", err)
			continue
		}
		go s.broadcast(&r)
	}
}

func (s *Session) broadcast(m *Message) {

	for ws := range s.conns {
		if err := websocket.JSON.Send(ws.connInfo, m); err != nil {
			fmt.Println("error writing to client:", err)
		}
	}
}

// for {
// 	var m Message
// 	if err := websocket.JSON.Receive(ws, &m); err != nil {
// 		log.Println(err)
// 		break
// 	}
// 	log.Println("Received message:", m.Message)
// 	// send a response
// 	m2 := Message{"Thanks for the message!"}
// 	if err := websocket.JSON.Send(ws, m2); err != nil {
// 		log.Println(err)
// 		break
// 	}
// }

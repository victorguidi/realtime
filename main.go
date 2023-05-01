package main

import (
	"log"
	"net/http"

	"golang.org/x/net/websocket"
)

// TODO: add a mutex to protect the map
// TODO: find a way to add TLS support (Done?)
// TODO: implement sessions

func main() {

	store, err := NewDatabase("./databases/database.db")
	if err != nil {
		log.Fatal(err)
	}
	store.Init()

	server := NewServer(":8080", store)

	// // http.Handle("/ws/orderbook", websocket.Handler(server.handleWSOrderbook))
	// // http.Handle("/wss/auth/orderbook", websocket.Handler(server.handleWSOrderbookWithAuth))

	http.HandleFunc("/api/getSessions", server.handleGetOpenSessions)
	http.HandleFunc("/api/getUsers", server.handleGetUsers)
	http.HandleFunc("/api/registerUser", server.handleRegisterNewUser)
	http.HandleFunc("/api/registerSession", server.handleCreateNewSession)
	http.HandleFunc("/api/registerUserToSession", server.handleUserSessions)
	http.HandleFunc("/api/login", server.handleLoginUser)
	http.Handle("/wss/login", websocket.Handler(server.handleWs))
	log.Fatal(http.ListenAndServeTLS(server.listenAddr, "./selfCertificate/server.crt", "./selfCertificate/server.key", nil))
	// log.Fatal(http.ListenAndServe(server.listenAddr, nil))
}

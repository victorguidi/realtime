package main

import (
	"log"
	"net/http"
)

// TODO: add a mutex to protect the map
// TODO: find a way to add TLS support (Done?)
// TODO: implement sessions
// TODO: Implement the api server, in order to login the users and save the sessions
// TODO: Implemnt db

func main() {

	store, err := NewDatabase("./databases/database.db")
	if err != nil {
		log.Fatal(err)
	}
	store.Init()

	server := NewServer(":8080", store)
	//
	// // http.Handle("/ws/orderbook", websocket.Handler(server.handleWSOrderbook))
	// // http.Handle("/wss/auth/orderbook", websocket.Handler(server.handleWSOrderbookWithAuth))
	http.HandleFunc("/api/getSessions", server.handleGetOpenSessions)
	// http.HandleFunc("/api/registerUser", server.handleRegisterNewUser)
	// http.HandleFunc("/api/registerSession", server.handleCreateNewSession)
	// http.HandleFunc("/api/registerUserToSession", server.handleAddUsersToSession)
	// http.HandleFunc("/wss/login", server.handleUpgradeToWsSession)
	log.Fatal(http.ListenAndServeTLS(server.listenAddr, "./selfCertificate/server.crt", "./selfCertificate/server.key", nil))
}

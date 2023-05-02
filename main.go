package main

import (
	"log"
	"mime"
	"net/http"

	"golang.org/x/net/websocket"
)

// TODO: add a mutex to protect the map

func main() {

	store, err := NewDatabase("./databases/database.db")
	if err != nil {
		log.Fatal(err)
	}
	store.Init()

	server := NewServer(":8080", store)

	mime.AddExtensionType(".js", "application/javascript")
	http.Handle("/", http.FileServer(http.Dir("static")))
	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "static/login.html")
	})
	http.HandleFunc("/api/getSessions", withJWTAuth(server.handleGetOpenSessions))
	http.HandleFunc("/api/getUsers", server.handleGetUsers)
	http.HandleFunc("/api/registerUser", server.handleRegisterNewUser)
	http.HandleFunc("/api/registerSession", server.handleCreateNewSession)
	http.HandleFunc("/api/registerUserToSession", server.handleUserSessions)
	http.HandleFunc("/api/login", server.handleLoginUser)
	http.Handle("/wss/login", websocket.Handler(server.handleWs))
	log.Fatal(http.ListenAndServeTLS(server.listenAddr, "./selfCertificate/server.crt", "./selfCertificate/server.key", nil))
}

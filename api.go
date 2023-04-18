package main

import "fmt"

// TODO: Sessions should be stored as raw data in a Database
// TODO: Users should be stored as raw data in a Database

type ApiSession struct {
	Sessions []Server
	Users    []User
}

type Api struct {
	listenAddr string
	*ApiSession
}

func NewApi(listenAddr string) {

	a := Api{
		listenAddr: "8080",
		ApiSession: &ApiSession{
			Sessions: make([]Server, 0),
			Users:    make([]User, 0),
		},
	}

	fmt.Println(a)
}

func (a *Api) Start() {}

func (a *Api) handleLoginUser() {}

func (a *Api) handleCreateSession() {}

func (a *Api) handleProxySession() {}

func (a *Api) handleGetSessions() {}

func (a *Api) handleUpdateSessions() {}

func (a *Api) handleDeleteSessions() {}

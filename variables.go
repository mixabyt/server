package main

import (
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

var (
	countClients = 0
	id           = 0
	RoomsID      = 0

	clients     = &Clients{clientsmap: make(map[int]*Client)}
	queueUsers  = &QueueUsers{Queue: make([]*Client, 0, 1)}
	rooms       = &Rooms{Rooms: make(map[int]*Room)}
	subMainMenu = make(map[int]*Client) // підписка на лічильник користувачів

	mu sync.Mutex
)

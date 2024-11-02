package main

import (
	"time"

	"github.com/gorilla/websocket"
)

type Clients struct {
	clientsmap map[int]*Client
}

type Client struct {
	ID           int
	RoomID       int
	Conn         *websocket.Conn
	LastActivity time.Time
}

func (c *Clients) AddNewUser(clientID int, conn *websocket.Conn) {
	mu.Lock()
	client := &Client{ID: clientID, Conn: conn, LastActivity: time.Now().Add(10 * time.Second), RoomID: -1}
	c.clientsmap[clientID] = client
	subMainMenu[clientID] = client
	countClients++
	id++
	mu.Unlock()
}

func (c *Clients) DeleteUser(clientID int) {
	mu.Lock()
	delete(c.clientsmap, clientID)
	countClients--
	mu.Unlock()
}

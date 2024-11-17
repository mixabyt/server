package main

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

type Rooms struct {
	Rooms map[int]*Room
}

type Room struct {
	Clients map[int]*Client
}

func (r *Rooms) CreateRoom(client *Client) {
	mu.Lock()
	defer mu.Unlock()

	newRoom := &Room{Clients: make(map[int]*Client)}
	newRoom.Clients[client.ID] = client
	r.Rooms[RoomsID] = newRoom
	client.RoomID = RoomsID

	log.Printf("id[%d] create and add to room[%d]\n", client.ID, RoomsID)
	RoomsID++
}

func (r *Rooms) AddToRoom(roomID int, client *Client) {
	mu.Lock()
	defer mu.Unlock()

	client.RoomID = roomID
	room, ok := r.Rooms[roomID]
	if !ok {

		log.Printf("roomid: %d not ok; clientroomid: %d, clientid: %d", roomID, client.RoomID, client.ID)
	}

	room.Clients[client.ID] = client
	log.Printf("id[%d] add to room[%d]\n", client.ID, roomID)
}

func (r *Rooms) DeleteRoom(roomID int) {
	log.Printf("delete room[%d]", roomID)
	delete(r.Rooms, roomID)
}

func (r *Room) SendMessage(client *Client, message interface{}) {
	text, _ := json.Marshal(message)
	mu.Lock()
	defer mu.Unlock()
	for _, v := range r.Clients {
		if v != client {
			log.Printf("message from[%d] to room[%d]", client.ID, client.RoomID)
			v.Conn.WriteMessage(websocket.TextMessage, text)

		}
	}
}

func (r *Room) KickFromRoom(client *Client) (int, int) {
	// (room id, how much user there is)
	rID := client.RoomID
	log.Printf("user[%d] has been kicked from room[%d]", client.ID, client.RoomID)
	client.RoomID = -1
	r.SendMessage(client, &DeleteNotice{Type: "roomDeletionNotice"})
	delete(r.Clients, client.ID)
	return rID, len(r.Clients)
}

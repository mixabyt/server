package main

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

func OnDisconect(client *Client, done chan<- bool) {
	// сповістити користувача по кімнаті що інший покинув
	if client.RoomID != -1 {
		rID, roomlenght := rooms.Rooms[client.RoomID].KickFromRoom(client)
		if roomlenght == 0 {
			rooms.DeleteRoom(rID)
		}
	}
	clients.DeleteUser(client.ID)
	queueUsers.DeleteFromQueue()
	OncounterNotify()
	done <- true
}

func OncounterNotify() {
	mu.Lock()
	for _, conn := range subMainMenu {

		sMM := &UpdateCountUser{Type: "subMainMenu", Count: countClients}
		data, _ := json.Marshal(sMM)
		conn.Conn.WriteMessage(websocket.TextMessage, data)
	}
	mu.Unlock()
}

func onSubMainMenu(client *Client, message []byte) {
	subuser := &SubMain{}
	json.Unmarshal(message, subuser)
	if subuser.Subscription {
		mu.Lock()
		subMainMenu[client.ID] = client
		mu.Unlock()
	} else {
		mu.Lock()
		delete(subMainMenu, client.ID)
		mu.Unlock()
	}
}

func onFindInterlocutor(client *Client) {
	interlocutor, inqueue := queueUsers.AddtoQueue(client)
	if !inqueue {
		log.Printf("rooms.CreateRoom(client)")
		rooms.CreateRoom(client)
	} else {
		log.Printf("перед queueUsers.DeleteFromQueue()")
		queueUsers.DeleteFromQueue()
		rooms.AddToRoom(interlocutor.RoomID, client)
		data, _ := json.Marshal(&FindInterlocutor{Type: "findInterlocutor"})
		for _, conn := range rooms.Rooms[interlocutor.RoomID].Clients {
			conn.Conn.WriteMessage(websocket.TextMessage, data)
		}
	}
}

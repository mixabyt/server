package main

import (
	"encoding/json"
	"fmt"
	"log"
	"sync"
	"time"
)

func ListenClient(clientID int, done chan<- bool, wg *sync.WaitGroup) {
	defer wg.Done()
	mu.Lock()
	client := clients.clientsmap[clientID]
	mu.Unlock()
	for {
		_, message, err := client.Conn.ReadMessage()
		if err != nil {
			log.Printf("user disconected: %s[%d]", client.Conn.RemoteAddr(), clientID)
			OnDisconect(client, done)
			break
		}

		typemessage := &BaseMessage{}
		err = json.Unmarshal(message, typemessage)
		if err != nil {
			log.Printf("Failed to unmarshal JSON: %s", err)
			continue
		}

		switch typemessage.Type {

		case "heartbeat":
			client.LastActivity = time.Now().Add(10 * time.Second)
		case "subMainMenu":
			onSubMainMenu(client, message)
		case "findInterlocutor":
			onFindInterlocutor(client)
		case "stopFindingInterlocutor":
			rID, _ := rooms.Rooms[client.RoomID].KickFromRoom(client)
			rooms.DeleteRoom(rID)
			queueUsers.DeleteFromQueue()
		case "message":
			data := &TextMessage{}
			json.Unmarshal(message, data)
			rooms.Rooms[client.RoomID].SendMessage(client, data)

		case "leaveDialog":
			rID, roomlenght := rooms.Rooms[client.RoomID].KickFromRoom(client)
			if roomlenght == 0 {
				rooms.DeleteRoom(rID)
			}
		default:
			fmt.Printf("невідомий тип повідомлення: %s", typemessage.Type)
		}

	}
}

package main

import (
	"encoding/json"
	"log"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

func OnHeartbeat(clientID int, done <-chan bool, wg *sync.WaitGroup) {
	defer wg.Done()

	for {
		select {
		case <-time.After(10 * time.Second):
			// log.Println("send heartbeat to:", clientID)
			mu.Lock()
			client, exists := clients.clientsmap[clientID]
			if !exists {
				mu.Unlock()
				return
			}
			// log.Printf("client[%d] last activity:%d", clientID, time.Now().Second()-client.LastActivity.Second())
			if time.Since(client.LastActivity) > 10*time.Second {

				log.Printf("client[%d] didn't respond to heartbeat, disconnecting", clientID)
				client.Conn.Close()
				mu.Unlock()
				return
			}
			mu.Unlock()

			// Відправляємо heartbeat повідомлення
			heartbeat := &BaseMessage{Type: "heartbeat"}
			data, _ := json.Marshal(heartbeat)
			err := client.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Printf("Failed to send heartbeat to client %d: %s", clientID, err)
				return
			}
		case <-done:
			return

		}
	}
}

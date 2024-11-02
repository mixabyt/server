package main

import (
	"log"
	"net/http"
	"sync"
)

func HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("fail upgrade to web socket: %s", err)
	}
	defer conn.Close()

	mu.Lock()
	clientID := id
	mu.Unlock()

	clients.AddNewUser(clientID, conn)
	log.Printf("user connected: %s [%d]", conn.RemoteAddr(), clientID)
	OncounterNotify()

	var wg sync.WaitGroup
	wg.Add(1)
	done := make(chan bool) // змінна для контролю горутини OnHeartbeat (щоб завершини коли користувач від'єднається)
	go OnHeartbeat(clientID, done, &wg)

	wg.Add(1)
	go ListenClient(clientID, done, &wg)
	wg.Wait()

}

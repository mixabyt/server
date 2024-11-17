package main

import "log"

type QueueUsers struct {
	Queue []*Client
}

func (q *QueueUsers) AddtoQueue(c *Client) (*Client, bool) {
	mu.Lock()
	defer mu.Unlock()
	if len(q.Queue) == 0 {
		log.Printf("черга пуста")
		q.Queue = append(q.Queue, c)
		return nil, false
	} else {
		log.Printf("черга не пуста client: %s", q.Queue[0].Conn.RemoteAddr())
		return q.Queue[0], true
	}
}

func (q *QueueUsers) DeleteFromQueue() {
	mu.Lock()
	q.Queue = q.Queue[:0]
	mu.Unlock()
}

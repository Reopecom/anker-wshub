package main

import (
	"encoding/json"
	"log"
	"sync"
)

type Subscription struct {
	mu          sync.RWMutex
	subscribers map[string][]*Client
	closed      bool
}

func NewSubscription() *Subscription {
	ps := &Subscription{}
	ps.subscribers = make(map[string][]*Client)
	return ps
}

func (ps *Subscription) Subscribe(topic string, client *Client) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	for _, c := range ps.subscribers[topic] {
		if c == client {
			// already subscribed
			log.Println("already subscribed")
			return
		}
	}

	ps.subscribers[topic] = append(ps.subscribers[topic], client)

}

func (ps *Subscription) UnSubscribe(topic string, client *Client) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	for i, c := range ps.subscribers[topic] {
		if c == client {
			ps.subscribers[topic][i] = ps.subscribers[topic][len(ps.subscribers[topic])-1]
			ps.subscribers[topic] = ps.subscribers[topic][:len(ps.subscribers[topic])-1]
			return
		}
	}

}

func (ps *Subscription) Publish(topic string, msg interface{}) {
	ps.mu.RLock()
	defer ps.mu.RUnlock()

	if ps.closed {
		return
	}
	data, err := json.Marshal(msg)
	if err != nil {
		log.Printf("error marshalling message; %s\n", err)
		return
	}
	for _, ch := range ps.subscribers[topic] {
		go func(c *Client) {
			c.send <- data
		}(ch)
	}
}

func (ps *Subscription) Close() {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	if !ps.closed {
		ps.closed = true
	}
}

func (ps *Subscription) Status() interface{} {
	ps.mu.RLock()
	defer ps.mu.RUnlock()
	data := make(map[string]interface{})

	for topic, clients := range ps.subscribers {
		data[topic] = len(clients)
	}
	return data

}

func (ps *Subscription) RemoveClient(client *Client) {
	ps.mu.Lock()
	defer ps.mu.Unlock()

	for topic, c := range ps.subscribers {
		for i, client_ := range c {
			if client_ == client {
				ps.subscribers[topic][i] = ps.subscribers[topic][len(ps.subscribers[topic])-1]
				ps.subscribers[topic] = ps.subscribers[topic][:len(ps.subscribers[topic])-1]
			}
		}
	}

}

package main

import (
	"encoding/json"
	"log"
	"sync"
)

type Message struct {
	ClientID string `json:"clientID"`
	Text     string `json:"text"`
	MsgType  string `json:"type"`
}

type Hub struct {
	sync.RWMutex
	// track the registered client
	clients map[*Client]bool

	// register channel
	register chan *Client

	// unregister channel
	unregister chan *Client

	// broadcast channel - msg from client
	broadcast chan *Message

	//messages history
	messages []*Message
}

func newHub() *Hub {
	log.Printf("Creating Room")
	return &Hub{
		clients:    make(map[*Client]bool),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		broadcast:  make(chan *Message),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			// listening for client register
			{
				// check if client exists and add the client entry
				h.Lock()
				log.Printf("client registered %s", client.id)
				h.clients[client] = true
				h.Unlock()

				// send history messages to client send channel
				for _, v := range h.messages {
					msgData, _ := json.Marshal(v)
					client.send <- []byte(string(msgData))
				}
			}
		case client := <-h.unregister:
			// listening for client unregister
			{
				// check if client exists and remove the client entry
				h.Lock()
				if _, clientExist := h.clients[client]; clientExist {
					h.clients[client] = false
					close(client.send)
					log.Printf("client unregistered %s", client.id)
					// h.messages = append(h.messages, logMsg)
					delete(h.clients, client)
				}
				h.Unlock()
			}
		case m := <-h.broadcast:
			// listening for msg to broadcast
			{
				// send message to all the active/registered clients
				h.RLock()
				log.Printf("%+v\n", m)
				h.messages = append(h.messages, m)
				msgData, _ := json.Marshal(m)
				for c := range h.clients {
					select {
					case c.send <- []byte(string(msgData)):
					default:
						close(c.send)
						delete(h.clients, c)
					}
				}
				h.RUnlock()
			}
		}
	}
}

type HubList struct {
	Hubs map[string]*Hub
}

func newHubList() *HubList {
	return &HubList{
		Hubs: make(map[string]*Hub),
	}
}

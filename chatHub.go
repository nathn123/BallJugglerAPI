package main

import "fmt"

//Chat struct to handle chat for game
type Chat struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

//NewHub creates a new chat hub for a game
func NewHub() *Chat {
	return &Chat{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[*Client]bool),
	}

}

func (h *Chat) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			fmt.Println("CLIENT REGISTERED")
		case client := <-h.unregister:
			fmt.Println("CLIENT UNREGISTERED")
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
			}
		case message := <-h.broadcast:
			fmt.Println("CLIENT broadcasting")
			for client := range h.clients {
				select {
				case client.send <- message:
				default:
					fmt.Println("MESSAGE NOT SENT")
					delete(h.clients, client)
					close(client.send)

				}
			}

		}
	}
}

// TODO msg routing for chat window / server requests
// TODO add close

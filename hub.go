package main

import "fmt"

/*
Hub manage clients and message
*/
type Hub struct {
	app        *App
	black      *Client
	white      *Client
	fromClient chan []byte
	register   chan *Client
	unregister chan *Client
}

func newHub(app *App) *Hub {
	return &Hub{
		app:        app,
		fromClient: make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			if h.black == nil {
				fmt.Println("register black")
				h.black = client
				message := h.app.Broadcast()
				h.black.send <- message
			} else if h.white == nil {
				fmt.Println("register white")
				h.white = client
				message := h.app.Broadcast()
				h.white.send <- message
			} else {
				fmt.Println("no register")
			}
		case client := <-h.unregister:
			close(client.send)
		case data := <-h.fromClient:
			h.app.Read(data)
			message := h.app.Broadcast()
			h.black.send <- message
			h.white.send <- message
			// case data := <-h.fromServer:
		}
	}
}

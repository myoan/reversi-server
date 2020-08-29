package main

/*
Hub manage clients and message
*/
type Hub struct {
	event      GameEvent
	fromClient chan *Data
	register   chan *Client
	unregister chan *Client
}

func newHub(event GameEvent) *Hub {
	return &Hub{
		event:      event,
		fromClient: make(chan *Data),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.event.OnRegister(client)
		case client := <-h.unregister:
			h.event.OnUnregister(client)
			close(client.send)
		case data := <-h.fromClient:
			h.event.OnFromClient(data.client, data.data)
		}
	}
}

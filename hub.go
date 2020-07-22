package main

// Hub connection hub for managing clients
type Hub struct {
	clients                   map[string]*Client // map of connected clients
	rooms                     map[string]*Room   // map of all rooms
	broadcast                 chan []byte        // channel to broadcast messages to all clients
	add, remove               chan *Client       // channel for adding clients
	addToRoom, removeFromRoom chan *RoomEntry    // channel for removing clients
	broadcastToRoom           chan *RoomMessage  // broadcast message to room
	createRoom, deleteRoom    chan *Room         // channel for room creation deletion
}

// create new instance of hub
func NewHub() *Hub {
	hub := &Hub{
		clients:        make(map[string]*Client),
		broadcast:      make(chan []byte),
		add:            make(chan *Client),
		addToRoom:      make(chan *RoomEntry),
		remove:         make(chan *Client),
		removeFromRoom: make(chan *RoomEntry),
		createRoom:     make(chan *Room),
		deleteRoom:     make(chan *Room),
	}
	go hub.Run()
	return hub
}

// run hub and manage clients
func (h *Hub) Run() {
	for {
		select {
		// add new client and update counter
		case client := <-h.add:
			h.clients[client.id] = client
		case entry := <-h.addToRoom:
			if h.rooms[entry.room] != nil {
				h.rooms[entry.room].clients[entry.client.id] = entry.client
			}
		case client := <-h.remove:
			if h.clients[client.id] != nil {
				delete(h.clients, client.id)
			}
		case entry := <-h.removeFromRoom:
			if h.rooms[entry.room] != nil && h.rooms[entry.room].clients[entry.client.id] != nil {
				delete(h.rooms[entry.room].clients, entry.client.id)
			}

		case room := <-h.createRoom:
			h.rooms[room.id] = room

		case room := <-h.deleteRoom:
			delete(h.rooms, room.id)
		case message := <-h.broadcastToRoom:
			if h.rooms[message.room] == nil {
				break
			}
			for id, client := range h.rooms[message.room].clients {
				select {
				// push message to client send channel
				case client.send <- message.message:
					// close channel when buffer is full
					// delete client and update counter
				default:
					if h.clients[id] != nil {
						close(client.send)
						delete(h.clients, id)
						delete(h.rooms[message.room].clients, id)
					}
				}
			}
			// broadcast to all clients
		case message := <-h.broadcast:
			for id, client := range h.clients {
				select {
				// push message to client send channel
				case client.send <- message:
					// close channel when buffer is full
					// delete client and update counter
				default:
					if h.clients[id] != nil {
						close(client.send)
						delete(h.clients, id)
					}
				}
			}
		}
	}
}

package main

import (
	"github.com/google/uuid"
)

type Room struct {
	id      string             // room id
	owner   string             // owner id  of the room
	clients map[string]*Client // map of clients in the room
}

type RoomEntry struct {
	room   string
	client *Client
}

type RoomMessage struct {
	room    string
	message []byte
}

func NewRoom(owner string) *Room {
	room := &Room{
		id:      uuid.New().String(),
		owner:   owner,
		clients: make(map[string]*Client),
	}
	return room
}

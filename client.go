package main

import (
	"log"

	"github.com/gofiber/websocket"
	"github.com/google/uuid"
)

// Client for handling read/write
type Client struct {
	id   string          // a client id
	conn *websocket.Conn // fiber/fasthttp websocket connnection
	hub  *Hub            // reference to hub
	send chan []byte     // channel to recevice messages from hub
}

// Create new Client
func NewClient(h *Hub, c *websocket.Conn) {
	// create a new client and push to hub
	client := &Client{
		id:   uuid.New().String(),
		conn: c,
		hub:  h,
		send: make(chan []byte, 256),
	}
	client.hub.add <- client
	// listen for writes in goroutine
	go client.write()
	client.read()
}

// read messages from websocket
func (c *Client) read() {
	defer func() {
		// remove client from hub and close connection once we are done
		c.hub.remove <- c
		if c.conn.Conn != nil {
			_ = c.conn.Close()
		}
	}()
	for {
		// read messages
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			break
		}
		// print message
		log.Println("Message", string(message))
		// send to write ( like echo )
		c.send <- message
	}
}

// write messages to websocket
func (c *Client) write() {
	defer func() {
		// remove client from hub and close connection once we are done
		c.hub.remove <- c
		if c.conn.Conn != nil {
			_ = c.conn.Close()
		}
	}()
	for message := range c.send {
		// send current message from channel
		err := c.conn.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			return
		}
		// send all others from channel buffer
		n := len(c.send)
		for i := 0; i < n; i++ {
			err = c.conn.WriteMessage(websocket.TextMessage, <-c.send)
			if err != nil {
				return
			}
		}
	}
}

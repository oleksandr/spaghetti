package spaghetti

import (
	"golang.org/x/net/websocket"
	"log"
)

// Handle connection from clients for publications only
func PubWebHandler(ws *websocket.Conn) {
	webHandler(ws, ConnectionTypePub)
}

// Handle connection from clients for subscriptions only
func SubWebHandler(ws *websocket.Conn) {
	webHandler(ws, ConnectionTypeSub)
}

// Handle connection from clients for both publications and subscriptions
func PubSubWebHandler(ws *websocket.Conn) {
	webHandler(ws, ConnectionTypePubSub)
}

func webHandler(ws *websocket.Conn, connType ConnectionType) {
	log.Printf("%vWebHandler: New connection from %v", connType, ws.RemoteAddr())

	// Create a new connection structure
	c := NewConnection(ws)
	switch connType {
	case ConnectionTypePub:
		c.Type = ConnectionTypePub
	case ConnectionTypeSub:
		c.Type = ConnectionTypeSub
	default:
		c.Type = ConnectionTypePubSub
	}

	// Notify hub to register the new connection
	DefaultHub.Register <- c

	// Defer the hub notification to unregister when this connection is closed
	defer func() {
		log.Printf("PubWebHandler: deferred call to remove connection from registry %v", ws.RemoteAddr())
		DefaultHub.Unregister <- c
	}()

	// Start writing in a separate goroutine
	go c.Writer()

	// Start reading from socket in the current goroutine
	c.Reader()
}

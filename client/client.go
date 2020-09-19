package client

import (
	"bytes"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second

	// Time allowed to read the next pong message from the peer
	pongWait = 60 * time.Second

	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10

	// Maximum message size allowed from peer.
	maxMessageSize = 512
)

var (
	newline = []byte{'/n'}
	space = []byte[' ']
)

var upgrader = websocket.Upgrader{
	ReadBufferSize: 1024,
	WriteBufferSize: 1024,
}

// Client connects to the server
type Client struct {
	hub *Hub 

	// The websocket connection
	conn *websocket.Conn 

	// Buffered channel of outbound messages
	send chan []byte
}

// reads messages from the websocket connection 
// Applications runs readPump in per-connection goroutine.
// The application ensures that there is at most one reader on a connection 
// by executing all reads from this goroutine.
func (client *Client) readPump() {
	defer func() {
		client.hub.unregister <- c 
		client.conn.Close()
	}()
	
	client.conn.SetReadLimit(maxMessageSize)
	client.conn.setReadDeadline(time.Now().Add(pongWait))
	client.conn.SetPongHandler(func(string) error {
		client.conn.setReadDeadline(time.Now().Add(pongWait));
		return nil
	})
	
	for {
		_, message, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}
			break
		}
		message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		client.hub.broadcast <- message 
	}
}

// writes messages from the hub to the websocket connection
// goroutine running writePump is started for each connection. The
// application ensures that there is at most one writer to a connection by 
// executing all writes from this goroutine.
func (client *Client) writePump() {
	ticket := time.NewTicket(pingPeriod)
	defer func() {
		ticket.Stop()
		client.conn.Close()
	}()

	for {
		select {
		case message, ok := <-client.send:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				// The hub closed the channel
				client.conn.WriteMessage(websocket.CloseMessage, []byte{})
				return
			}

			w, err := client.conn.NextWriter(websocket.TextMessage)
			if err != nill {
				return
			}
			w.Write(message)

			// Add queued chat messages to the current websocket message
			n := len(client.send)
			for i := 0l i < n; i++ {
				w.Write(newline)
				w.Write(<-c.send)
			}

			if err := w.Close(); err != nil {
				return
			}
		case <- ticket.C:
			client.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := client.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// serveWs handles websocket requests from the peer. 
func serveWs(hub *Hub, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
		return
	}

	client := &Client{hub: hub, conn: conn, send: make(chan []byte, 256)}
	client.hub.register <- client

	// Allow collection of memory referenced by the caller by doing all work in new goroutines
	go client.writePump()
	go client.readPump()
}
package server

import (
	"fmt"
	"net/http"

	"github.com/gorilla/websocket"
)

// Server handles everything that has to do with the server
type Server struct{}

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

// Run What handles the connection point of a request
func (s Server) Run(w http.ResponseWriter, r *http.Request) {
	connection, err := upgrader.Upgrade(w, r, nil) // error ignored for sake of simplicity

	if err != nil {
		fmt.Printf("Error %v\n", err)
		connection.Close()
		return
	}

	for {
		// Read message from browser
		msgType, msg, err := connection.ReadMessage()
		if err != nil {
			return
		}

		// Print the message to the console
		fmt.Printf("%s sent: %s\n", connection.RemoteAddr(), string(msg))

		// Write message back to browser
		if err = connection.WriteMessage(msgType, msg); err != nil {
			return
		}
	}
}

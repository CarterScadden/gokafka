package main

import (
	"gokafka/server"
	"net/http"
)

func main() {
	server := server.Server{}

	http.HandleFunc("/socket", server.Run)

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "websockets.html")
	})

	http.ListenAndServe(":8080", nil)
}

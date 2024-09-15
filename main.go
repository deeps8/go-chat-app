package main

import (
	"errors"
	"log"
	"net/http"
)

var port = ":42069"

func ServeHTML(w http.ResponseWriter, r *http.Request) {

}

func main() {
	mux := http.NewServeMux()

	// hub := newHub()
	// go hub.run()

	hubList := newHubList()

	mux.HandleFunc("GET /", func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, "index.html")
	})

	// mux.HandleFunc("GET /ws", func(w http.ResponseWriter, r *http.Request) {
	// 	// handle the client websocket request
	// 	serveWS(hub, w, r)
	// })

	// for rooms
	mux.HandleFunc("GET /ws", func(w http.ResponseWriter, r *http.Request) {
		// handle the client websocket request
		serveRoomWS(hubList, w, r)
	})

	serverErr := http.ListenAndServe(port, mux)
	if errors.Is(serverErr, http.ErrServerClosed) {
		log.Fatalf("Something went wrong, Server closed on %v", port)
	} else {
		log.Printf("Server started at %v", port)
	}
}

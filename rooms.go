package main

import (
	"fmt"
	"log"
	"net/http"
)

func serveRoomWS(hubList *HubList, w http.ResponseWriter, r *http.Request) {
	// roomId and user from params
	log.Printf("url %v", r.URL.String())
	roomid := r.URL.Query().Get("roomid")
	userid := r.URL.Query().Get("userid")
	if roomid == "" || userid == "" {
		log.Printf("Roomid: %v\nUserid: %v", roomid, userid)
		log.Fatal("Room or User ids are invalid")
		return
	}
	// upgrading http connection to websocket
	conn, wsErr := connUpgrader.Upgrade(w, r, nil)
	if wsErr != nil {
		log.Printf("Error while upgrading http connection to websocket")
	}

	/*
		1. if roomid exists or not
		2. if not create one and register the client
		3. if yes register the client in that room
	*/
	hub, ok := hubList.Hubs[roomid]
	log.Printf("Roomid : %v", hub)
	if !ok {
		log.Printf("Room Does not exists")
		// room does not exists, create one
		hub = newHub()
		hubList.Hubs[roomid] = hub
		go hub.run()
	}

	// creating a new client for each time use joins
	newClient := &Client{id: userid, hub: hub, conn: conn, send: make(chan []byte, 512)}

	// registering the newly created client to hub
	newClient.hub.register <- newClient

	// now start the goroutines for read and write the data
	go newClient.writePump()
	go newClient.readPump()

	logMsg := &Message{ClientID: newClient.id, Text: fmt.Sprintf("User Joined %s", newClient.id), MsgType: "info"}
	newClient.hub.broadcast <- logMsg
}

package main

import (
	"fmt"
	"go-chat/internal"
	"log"
	"net/http"
)

func main() {
	fmt.Printf("Websocket listening on localhost:8080\n")
	ws := internal.NewWebsocketChat()
	http.HandleFunc("/chat", ws.UserConnectionHandler)
	go ws.UsersChatManager()
	log.Fatalln(http.ListenAndServe("localhost:8080", nil))
}

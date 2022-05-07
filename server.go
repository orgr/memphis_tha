package main

import (
	"fmt"
	"log"
	"net/http"
	"encoding/json"
)

const (
	route = "sendMessage"
	httpPort = "8080"
)

type message struct {
	Author string `json:"author"`
	Data string `json:"data"`
}

func main() {
	http.HandleFunc("/" + route, handler)
	log.Println("Starting HTTP server")
	log.Fatal(http.ListenAndServe(":" + httpPort, nil))
}

func handler(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "POST":
		var receivedMessage message
		err := json.NewDecoder(r.Body).Decode(&receivedMessage)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		log.Println(receivedMessage)

		fmt.Fprintf(w, "new message: %s\n", receivedMessage)
		publishMessage(receivedMessage)

	default:
		fmt.Fprintf(w, "Sorry, only Post method is supported\n")
	}
}

func publishMessage(message message) {

}

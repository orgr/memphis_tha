package main

import (
	"fmt"
	"log"
	"net/http"
	"encoding/json"

	"github.com/nats-io/nats.go"
)

const (
	route = "sendMessage"
	httpPort = "8080"
	natsPort = nats.DefaultURL
	subject = "messages"
)

type message struct {
	Author string `json:"author"`
	Data string `json:"data"`
}

var natsClient nats.Conn

func main() {
	http.HandleFunc("/" + route, handler)
	log.Println("Establising connection with NATS server")
	natsConnectionPtr, err := nats.Connect(natsPort)
	if err != nil {
		log.Fatal(err)
		return
	}
	natsClient = *natsConnectionPtr

	subscribeToNats(&natsClient)
	log.Println("Starting HTTP server")
	log.Fatal(http.ListenAndServe(":" + httpPort, nil))
}

func subscribeToNats(nc *nats.Conn) {
	nc.Subscribe(subject, func(msg *nats.Msg) {
        // Print message data
        data := string(msg.Data)
        fmt.Println(data)
    })
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
	natsClient.Publish(subject, []byte(message.Data))
}

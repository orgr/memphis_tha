package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/nats-io/nats-server/v2/server"
	"github.com/nats-io/nats.go"
)

const (
	route                  = "sendMessage"
	httpPort               = "8080"
	natsPort               = nats.DefaultURL
	serverReadinessTimeout = 4 * time.Second
	subject                = "messages"
)

type message struct {
	Author string `json:"author"`
	Data   string `json:"data"`
}

var natsServer *server.Server
var natsClient *nats.Conn
var jetStream nats.JetStreamContext

func main() {
	setupNatsServer()

	setupNatsClient()
	defer natsClient.Close()

	setupJetStream()

	setupStream()

	log.Println("Starting HTTP server")
	http.HandleFunc("/"+route, httpHandler)
	log.Fatal(http.ListenAndServe(":"+httpPort, nil))
}

func setupNatsServer() {
	log.Println("Start an embedded NATS server")

	var err error
	serverOptions := &server.Options{
		JetStream: true,
	}
	natsServer, err = server.NewServer(serverOptions)
	exitOnError(err)

	go natsServer.Start()
	if !natsServer.ReadyForConnections(serverReadinessTimeout) {
		log.Fatal("Not ready for connections")
	}
}

func setupNatsClient() {
	log.Println("Establising connection with NATS server")

	var err error
	natsClient, err = nats.Connect(natsPort)
	exitOnError(err)
}

func setupJetStream() {
	var err error
	jetStream, err = natsClient.JetStream()
	exitOnError(err)
}

func setupStream() {
	stream, err := jetStream.StreamInfo(subject)
	if err != nil {
		log.Print(err)
	}

	if stream != nil {
		log.Printf("Stream \"%s\" exists\n", subject)
		return
	}

	log.Printf("Creating stream \"%s\"\n", subject)
	_, err = jetStream.AddStream(&nats.StreamConfig{
		Name: subject,
	})
	exitOnError(err)
}

func httpHandler(w http.ResponseWriter, r *http.Request) {
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
	_, err := jetStream.Publish(subject, []byte(message.Data))
	if err != nil {
		log.Print(err)
	}
}

func exitOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

package main

import (
	"fmt"
	"log"
	// "net/http"
	// "time"
	// "encoding/json"

	"github.com/nats-io/nats.go"
)

const (
	route    = "sendMessage"
	natsPort = nats.DefaultURL
	subject  = "messages"
)

type message struct {
	Author string `json:"author"`
	Data   string `json:"data"`
}

var natsClient *nats.Conn
var jetStream nats.JetStreamContext

func main() {
	setupNatsClient()
	defer natsClient.Close()

	setupJetStream()

	flushMessages()

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

func flushMessages() {
	stream, err := jetStream.StreamInfo(subject)
	exitOnError(err)

	if stream == nil {
		log.Fatal("Stream \"%s\" doesn't exist\n", subject)
	}

	log.Printf("Consuming stream \"%s\"\n", subject)
	jetStream.Subscribe(subject, func(msg *nats.Msg) {
		data := string(msg.Data)
		fmt.Println(data)
	})

}

func exitOnError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

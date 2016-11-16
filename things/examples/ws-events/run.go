package main

import (
	"log"
	"github.com/zubairhamed/iot-suite-sdk-go.bak/things/client"
	"github.com/zubairhamed/iot-suite-sdk-go.bak/things/ws"
	"github.com/zubairhamed/iot-suite-sdk-go.bak/things/examples"
	"github.com/zubairhamed/iot-suite-sdk-go.bak/things"
)

func main() {
	cfg := &client.Configuration{
		SkipSslVerify: true,
		// Proxy: "http://localhost:3128",
	}

	conn, err := ws.Dial(
		examples.ENDPOINT_URL_WS,
		examples.USERNAME,
		examples.PASSWORD,
		examples.APITOKEN,
		cfg,
	)

	if err != nil {
		panic(err.Error())
	}

	obsEvents := make(chan *things.WSMessage)
	obsDelete := make(chan *things.WSMessage)
	obsCreate := make(chan *things.WSMessage)
	obsUpdate := make(chan *things.WSMessage)

	conn.ObserveEvents(obsEvents)
	log.Println("Start Observing all events")

	conn.ObserveCreatedEvents(obsCreate)
	log.Println("Start Observing created events")

	conn.ObserveDeletedEvents(obsDelete)
	log.Println("Start Observing deleted events")

	conn.ObserveUpdateEvents(obsUpdate)
	log.Println("Start Observing updated events")

	for {
		select {
		case obsMsg, _ := <-obsEvents:
			log.Println(">> Event", obsMsg.Topic)

		case obsMsg, _ := <-obsDelete:
			log.Println(">> Delete Event", obsMsg.Topic)

		case obsMsg, _ := <-obsCreate:
			log.Println(">> Create Event", obsMsg.Topic)

		case obsMsg, _ := <-obsUpdate:
			log.Println(">> Update Event", obsMsg.Topic)
		}
	}
}

package main

import (
	"github.com/zubairhamed/iot-suite-sdk-go/things/client"
	"github.com/zubairhamed/iot-suite-sdk-go/things/ws"
	"github.com/zubairhamed/iot-suite-sdk-go/things/examples"
	"github.com/zubairhamed/iot-suite-sdk-go/things"
	"time"
	"fmt"
)

func main() {
	cfg := &client.Configuration{
		SkipSslVerify: true,
		// Proxy: "http://localhost:3128",
	}

	fmt.Println("### Connecting to Things via WebSockets..")
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
	fmt.Println("### Connected.")

	obsEvents := make(chan *things.WSMessage)
	obsDelete := make(chan *things.WSMessage)
	obsCreate := make(chan *things.WSMessage)
	obsUpdate := make(chan *things.WSMessage)

	conn.ObserveEvents(obsEvents)
	fmt.Println("### Start Observing all events")

	conn.ObserveCreatedEvents(obsCreate)
	fmt.Println("### Start Observing created events")

	conn.ObserveDeletedEvents(obsDelete)
	fmt.Println("### Start Observing deleted events")

	conn.ObserveUpdateEvents(obsUpdate)
	fmt.Println("### Start Observing updated events")

	tickChan := time.NewTicker(time.Second * 60).C

	for {
		select {
		case obsMsg, _ := <-obsEvents:
			fmt.Println(">> Incoming Event       ", obsMsg.Topic)

		case obsMsg, _ := <-obsDelete:
			fmt.Println(">> Incoming Delete Event", obsMsg.Topic)

		case obsMsg, _ := <-obsCreate:
			fmt.Println(">> Incoming Create Event", obsMsg.Topic)

		case obsMsg, _ := <-obsUpdate:
			fmt.Println(">> Incoming Update Event", obsMsg.Topic)

		case <- tickChan:
			examples.PrintMemoryStats()
		}
	}
}


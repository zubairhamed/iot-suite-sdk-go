package main

import (
	"github.com/zubairhamed/paho.mqtt.golang"
	"time"
	"fmt"
	"os"
	"github.com/zubairhamed/iot-suite-sdk-go/things/examples"
	"github.com/zubairhamed/iot-suite-sdk-go/things/ws"
	"github.com/zubairhamed/iot-suite-sdk-go/things"
)

// Subscribes to an MQTT Broker, listening to a specific topic
// Payload of topic should be compliant to what's needed for IoT Things
var MQTT_BROKER = "tcp://iot.eclipse.org:1883"
var INCOMING_MQTT_TOPIC = "/zoob/iot/things/incoming"
var OUTGOING_MQTT_TOPIC = "/zoob/iot/things/outgoing"
var MQTT_CLIENTID = "things-mqtt-client"

func main() {
	mqttEvents := make(chan mqtt.Message)
	// Connect to an MQTT Broker, listening to a topic
	opts := mqtt.NewClientOptions().AddBroker(MQTT_BROKER).SetClientID(MQTT_CLIENTID)
	opts.SetKeepAlive(2 * time.Second)
	opts.SetPingTimeout(1 * time.Second)
	opts.SetDefaultPublishHandler(func(client mqtt.Client, msg mqtt.Message) {
		mqttEvents <- msg
	})
	mqttConn := mqtt.NewClient(opts)
	fmt.Println("Connecting to MQTT Broker", MQTT_BROKER)
	if token := mqttConn.Connect(); token.Wait() && token.Error() != nil {
		panic(token.Error())
	}

	fmt.Println("Subscribing to MQTT Topic", INCOMING_MQTT_TOPIC, "/#")
	if token := mqttConn.Subscribe(INCOMING_MQTT_TOPIC + "/#", 0, nil); token.Wait() && token.Error() != nil {
		fmt.Println(token.Error())
		os.Exit(1)
	}

	// Connect to IoT Things
	fmt.Println("Connecting to Bosch IoT Things", examples.ENDPOINT_URL_WS)
	thingsConn, err := ws.Dial(
		examples.ENDPOINT_URL_WS,
		examples.USERNAME,
		examples.PASSWORD,
		examples.APITOKEN,
		examples.DEFAULT_CLIENT_CONFIG,
	)
	if err != nil {
		panic(err.Error())
	}

	// Observe to IoT Things Events
	thingEvents := make(chan *things.WSMessage)
	fmt.Println("Observing changes from IoT Things..")
	thingsConn.ObserveEvents(thingEvents)

	for {
		select {
		// Handle incoming Things Events
		case thingsMsg, _ := <-thingEvents:
			t := thingsMsg.ValueAsThing()
			fmt.Println("[IOT THINGS >> ] Incoming Message", t.ThingId)

			// Create Topic
			topic := OUTGOING_MQTT_TOPIC + "/" + thingsMsg.Topic

			// Convert Things to JSON payload
			if content, err := thingsMsg.ValueAsThing().String(); err == nil {

				// Publish to MQTT
				fmt.Println(">>> Publishing to MQTT", topic)
				mqttConn.Publish(topic, 0, false, content)
			} else {
				fmt.Println("Error occured updating MQTT with payload from Things")
			}

		// Handle incoming MQTT Events
		case mqttMsg, _ := <- mqttEvents:
			topic := mqttMsg.Topic()
			payload := mqttMsg.Payload()
			fmt.Println("[MQTT >> ] Incoming Message", topic, payload)

			// Convert Payload to IoT Things Object
			nt := things.NewThingFromContent(payload)
			if nt.ThingId != "" {
				fmt.Println(">>> Publishing to IoT Things", nt.ThingId)
				if err := thingsConn.Update(nt); err != nil {
					fmt.Println("Error occured updating Things with payload from MQTT")
				}
			} else {
				fmt.Println("Invalid Thing Payload encountered from MQTT. Dropping message.")
			}

		}
	}
}

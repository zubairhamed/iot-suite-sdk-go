package main

import (
	"log"
	"github.com/satori/go.uuid"
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

	t := things.NewThing()
	t.ThingId = "com.zoob.mythicalthings:" + uuid.NewV4().String()
	t.Attributes["name"] = "NameAttribute"

	t, err = conn.Add(t)
	if err != nil || t == nil {
		panic(err.Error())
	}
	log.Println("Thing created. id:", t.ThingId)

	thingId := t.ThingId
	t, err = conn.Get(thingId)
	if err != nil || t == nil {
		panic(err.Error())
	}

	if t.ThingId != thingId {
		panic("Unequal thing ID returned")
	}
	log.Println("Got back thing. id:", t.ThingId)

	t.Attributes["prop"] = "val"
	err = conn.Update(t)
	if err != nil {
		panic(err.Error())
	}

	log.Println("Thing updated.")

	t, err = conn.Get(thingId)
	if err != nil || t == nil {
		panic(err.Error())
	}

	if t.Attributes["prop"] != "val" {
		log.Println(t)
		panic("Property 'prop' is not of value 'val'")
	}

	err = conn.Delete(thingId)
	if err != nil {
		panic(err.Error())
	}

	log.Println("Deleted thing")

	t, err = conn.Get(thingId)
	if t != nil {
		panic("Thing should have been deleted.")
	}

	log.Println("CRUD test completed")
}

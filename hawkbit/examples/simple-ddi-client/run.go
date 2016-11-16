package main

import (
	"github.com/zubairhamed/iot-suite-sdk-go/hawkbit/ddi"
	"github.com/zubairhamed/iot-suite-sdk-go/hawkbit/examples"
	"github.com/zubairhamed/iot-suite-sdk-go/hawkbit"
	"log"
	"fmt"
	"sync"
)

func main() {
	var TENANT_ID = "xx"
	var TARGET_ID = "xx"
	var TARGET_SECURITY_TOKEN = "xx"

	fmt.Println("Connecting to Hawkbit..")

	conn, err := ddi.Dial(example.HAWKBIT_HTTP_ENDPOINT,
		TENANT_ID,
		TARGET_ID,
		TARGET_SECURITY_TOKEN)

	if err != nil {
		fmt.Println("Uh oh, error occured connecting to HawkBit..")
		panic(err.Error())
	}

	updateEvents := make(chan *hawkbit.Message)

	fmt.Println("Waiting for updates..")
	conn.WaitForUpdate(updateEvents)

	for {
		select {
		case updateMsg, _ := <- updateEvent:
			handleUpdateEvent(conn, updateMsg)
		}
	}
}

func handleUpdateEvent(conn Connection, msg Message) {
	fmt.Println("An update was found. Let's update!")

	actionId := msg.GetActionId()

	// Tell the mothership that we're proceeding to update
	conn.UpdateActionStatus(actionId, hawkbit.STATUS_EXEC_PROCEEDING, hawkbit.STATUS_RESULT_NONE)

	// Get all downloadable chunks for this update
	chunks := msg.GetDownloadableChunks()

	artifacts := []*hawkbit.Artifact
	var wg sync.WaitGroup
	for _, c := range chunks {
		fmt.Println("Downloading chunk part", c.Part, ", name", c.Name, "and version", c.Version)
		wg.Add(len(c.Artifacts))
		for _, a := range c.Artifacts {
			downloadHref := a.Links.DownloadHttp.Href
			fileName := a.Filename

			// dereference pointer
			artifact := *a

			// Download artifact binaries
			go func() {
				defer wg.Done()
				log.Println("<<<<<< Downloading Artifact", fileName, " @ ", downloadHref)
				if client.DownloadArtifact(&artifact) == nil {
					log.Println("++++++ Downloaded HawkBit Artifact ", fileName, "of size", len(artifact.Content))

					// Save this artifact to list for use later
					artifacts = append(artifacts, artifact)
				} else {
					log.Println("Uh oh, we failed to download an artifact..", fileName, " @ ", downloadHref)
				}
			}()
		}
	}
	wg.Wait()
	fmt.Println("Download all completed.")
	startUpdating(artifacts)
}

func startUpdating(a []*hawkbit.Artifact) {
	fmt.Println("Got", len(a), "artifacts. Start doing some fancy updates with them ..")
}



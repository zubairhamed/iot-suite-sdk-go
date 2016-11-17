package main

import (
	"github.com/zubairhamed/iot-suite-sdk-go/hawkbit/ddi"
	"github.com/zubairhamed/iot-suite-sdk-go/hawkbit/examples"
	. "github.com/zubairhamed/iot-suite-sdk-go/hawkbit"
	"fmt"
	"sync"
	"time"
	"github.com/zubairhamed/iot-suite-sdk-go"
)

func main() {
	var TENANT_ID = ""
	var TARGET_ID = ""
	var TARGET_SECURITY_TOKEN = ""

	fmt.Println("Connecting to Hawkbit..")

	conn, err := ddi.Dial(example.HAWKBIT_HTTP_ENDPOINT,
		TENANT_ID,
		TARGET_ID,
		TARGET_SECURITY_TOKEN,
		nil,
	)

	if err != nil {
		fmt.Println("Uh oh, error occured connecting to HawkBit..")
		panic(err.Error())
	}

	updateEvents := make(chan *Message)

	poll := 3
	fmt.Println("Waiting for updates.. Polling set to", poll, "seconds")
	conn.WaitForUpdates(updateEvents, poll)

	// Channel for printing memory stats
	tickChan := time.NewTicker(time.Second * 60).C

	for {
		select {
		case updateMsg, _ := <- updateEvents:
			handleUpdateEvent(conn, updateMsg)

		case <- tickChan:
			common.PrintMemoryStats()
		}
	}
}

func handleUpdateEvent(conn *ddi.HawkbitDDIConnection, msg *Message) {
	fmt.Println("@@@@@@@@ An update was found. Let's update!")

	actionId := msg.ActionId

	// Tell the mothership that we're proceeding to update
	conn.UpdateActionStatus(actionId, STATUS_EXEC_PROCEEDING, STATUS_RESULT_NONE)

	// Create a channel to be notified which artifact is going to be downloaded
	downloadingCh := make(chan Artifact)

	// Create a channel to be notified of downloaded artifacts
	downloadedCh := make(chan Artifact)

	// Download all artifacts, specifying to download artifacts in parallel
	totalArtifacts, err := conn.DownloadArtifacts(downloadingCh, downloadedCh, true)
	if err != nil {
		// Something happened trying to download artifacts..
		// OMG! Exit in panic!!
		panic(err.Error())
	}

	var wg sync.WaitGroup
	wg.Add(totalArtifacts)

	// An array to store all our downloaded content
	downloadedArtifacts := []Artifact{}

	// Wait for notifications of artifact downloads, run this in a goroutine
	go func() {
		for {
			select {
			case artifact := <-downloadingCh:
				fmt.Println(">>>>>>>> Downloading HawkBit Artifact ", artifact.Filename, "from", artifact.Links.DownloadHttp.Href)

			case artifact := <-downloadedCh:
				// When an artifact has been completed, save its content to an array
				fmt.Println("<<<<<<<< Downloaded HawkBit Artifact  ", artifact.Filename, "of size", len(artifact.Content), "bytes")
				downloadedArtifacts = append(downloadedArtifacts, artifact)
				wg.Done()
			}
		}
	}()

	// Wait for completion of downloads
	wg.Wait()
	fmt.Println("@@@@@@@@ Download of artifacts completed.")

	// With all the content of artifacts we've collected, let's start an update
	startUpdating(conn, actionId, downloadedArtifacts)
}

func startUpdating(conn *ddi.HawkbitDDIConnection, actionId string, a []Artifact) {
	fmt.Println("@@@@@@@@ Got", len(a), "artifacts. Start doing some fancy updates with them ..")

	// Tell mothership that all's good and update was a success.
	// If update failed, use this instead:
	//	    client.UpdateActionStatus(actionId, hawkbit.STATUS_EXEC_CLOSED, hawkbit.STATUS_RESULT_FAILED)
	conn.UpdateActionStatus(actionId, STATUS_EXEC_CLOSED, STATUS_RESULT_SUCCESS)
	fmt.Println("@@@@@@@@ Update completed. Waiting for next update..")
}

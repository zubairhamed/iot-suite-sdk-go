package main

import (
	"github.com/zubairhamed/iot-suite-sdk-go/hawkbit/ddi"
	"github.com/zubairhamed/iot-suite-sdk-go/hawkbit/examples"
	. "github.com/zubairhamed/iot-suite-sdk-go/hawkbit"
	"fmt"
)

func main() {
	var TENANT_ID = "d683638a-8d35-49c5-bacd-2075f1216df8"
	var TARGET_ID = "ZubairMacbookPro"
	var TARGET_SECURITY_TOKEN = "PXZMtJQ05zp75HenVOPognki0D1WNBCa"

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

	fmt.Println("Waiting for updates..")
	conn.WaitForUpdates(updateEvents)

	for {
		select {
		case updateMsg, _ := <- updateEvents:
			handleUpdateEvent(conn, updateMsg)
		}
	}
}

func handleUpdateEvent(conn *ddi.HawkbitDDIConnection, msg *Message) {
	fmt.Println("An update was found. Let's update!")

	actionId := msg.ActionId

	// Tell the mothership that we're proceeding to update
	conn.UpdateActionStatus(actionId, STATUS_EXEC_PROCEEDING, STATUS_RESULT_NONE)

	// Create a channel to be notified of downloaded artifacts
	ch := make(chan *Artifact)

	// Download all artifacts, specifying to download artifacts in parallel
	conn.DownloadArtifacts(ch, true)

	// An array to store all our downloaded content
	downloadedArtifacts := []*Artifact{}

	// Wait till all downloads have completed.
	for {
		select {
			case artifact, closed := <-ch:
				if closed {
					fmt.Println("Download of artifacts completed.")
					// Channel was closed, meaning update has completed, so let's get outta here
					break
				} else {
					// When an artifact has been completed, save its content to an array
					fmt.Println("++++++ Downloaded HawkBit Artifact ", artifact.Filename, "of size", len(artifact.Content))
					downloadedArtifacts = append(downloadedArtifacts, artifact)
				}
		}
	}
	// With all the content of artifacts we've collected, let's start an update
	startUpdating(downloadedArtifacts)
}

func startUpdating(a []*Artifact) {
	fmt.Println("Got", len(a), "artifacts. Start doing some fancy updates with them ..")
}



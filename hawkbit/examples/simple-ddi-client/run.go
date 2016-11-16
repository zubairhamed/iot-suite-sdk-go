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

	client := ddi.NewDefaultHawkbitClient(example.HAWKBIT_HTTP_ENDPOINT,
		TENANT_ID,
		TARGET_ID,
		TARGET_SECURITY_TOKEN)

	client.OnUpdate(UpdateCallback)
	fmt.Println("Connecting to Hawkbit..")
	client.Start()
	fmt.Println("Waiting for updates..")
}

func UpdateCallback(client hawkbit.HawkbitClient, actionId string) {
	fmt.Println("@@@@@@ An update was found. Let's update!")

	// Tell the mothership that we're proceeding to update
	client.UpdateActionStatus(actionId, hawkbit.STATUS_EXEC_PROCEEDING, hawkbit.STATUS_RESULT_NONE)

	// Get all downloadable chunks for this update
	chunks := client.GetDownloadableChunks()

	// Do some parallel downloading
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
				} else {
					log.Println("Uh oh, we failed to download an artifact..", fileName, " @ ", downloadHref)
				}
			}()
		}
	}
	wg.Wait()

	fmt.Println("Download all completed.")
	startUpdating()
}

func startUpdating() {
	fmt.Println("Start doing some fancy updates..")
}


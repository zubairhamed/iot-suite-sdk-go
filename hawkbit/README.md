# iot-suite-hawkbit-go

### Direct Device Integration API
This uses HTTP to periodically poll HawkBit and to download updates.
It supports parallel download of artifacts and asychronous notification of artifact downloads via channels.

``` go
	// Create new connection to HawkBit
	conn, err := ddi.Dial(example.HAWKBIT_HTTP_ENDPOINT,
		TENANT_ID,
		TARGET_ID,
		TARGET_SECURITY_TOKEN,
		nil,
	)

	// Poll every 20 seconds
	pollSecs := 20

	// Start waiting for updates
	conn.WaitForUpdates(updateEvents, pollSecs)

	for {
		select {
		case updateMsg, _ := <- updateEvents:	// An update event has occured
			handleUpdateEvent(conn, updateMsg)
		}
	}

	func handleUpdateEvent(conn *ddi.HawkbitDDIConnection, msg *Message) {
		// Notify HawkBit that we're beginning to update the firmware
		conn.UpdateActionStatus(actionId, STATUS_EXEC_PROCEEDING, STATUS_RESULT_NONE)

		// Create two channels, one for being notified of artifact being downloaded and another to be notified
		// when an artifact has completed download
		downloadingCh := make(chan Artifact)
		downloadedCh := make(chan Artifact)
		totalArtifacts, err := conn.DownloadArtifacts(downloadingCh, downloadedCh, true)

		// Wait for all downloads to complete
		var wg sync.WaitGroup
		wg.Add(totalArtifacts)

		go func() {
			for {
				select {
				case artifact := <-downloadingCh:
					fmt.Println(">>>>>>>> Downloading HawkBit Artifact ", artifact.Filename, "from", artifact.Links.DownloadHttp.Href)

				case artifact := <-downloadedCh:
					fmt.Println("<<<<<<<< Downloaded HawkBit Artifact  ", artifact.Filename, "of size", len(artifact.Content), "bytes")
					wg.Done()
				}
			}
		}()

		// Wait for completion of downloads
		wg.Wait()

		// Downloads completed
	}
```

### Device Management Federation API

### Management API
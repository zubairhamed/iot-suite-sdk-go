package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"sync"
	"github.com/zubairhamed/iot-suite-sdk-go.bak/hawkbit/hawkbit"
	"github.com/zubairhamed/iot-suite-sdk-go.bak/hawkbit/hawkbit/ddi"
)

var PLAYER_PROCESS *exec.Cmd

func main() {
	log.Println("@@@@@@ Started Updater Daemon.")
	StartPlayer()
	client := ddi.NewDefaultHawkbitClient("https://rollouts-cs.apps.bosch-iot-cloud.com",
		"d683638a-8d35-49c5-bacd-2075f1216df8",
		"ZubairMacbookPro",
		"PXZMtJQ05zp75HenVOPognki0D1WNBCa")
	client.UseProxy("http://localhost:3128")
	client.OnUpdate(UpdateCallback)

	log.Println("@@@@@@ Waiting for updates.")
	client.Start()
}

func UpdateCallback(client hawkbit.HawkbitClient, actionId string) {
	log.Println("@@@@@@ An update was found. Let's update!")
	client.UpdateActionStatus(actionId, hawkbit.STATUS_EXEC_PROCEEDING, hawkbit.STATUS_RESULT_NONE)

	chunks := client.GetDownloadableChunks()

	// Do parallel downloading till completed
	var wg sync.WaitGroup
	for _, c := range chunks {
		chunkPart := c.Part
		chunkName := c.Name
		chunkVersion := c.Version

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
					// Save to Disk
					savePath := fmt.Sprintf("./hawkbit/%s/%s/%s", chunkPart, chunkName, chunkVersion)
					savedFile := fileName

					err := SaveContent(savePath, &artifact)
					if err != nil {
						panic(err.Error())
					}
					log.Println("++++++ Downloaded HawkBit Artifact ", savedFile)
				} else {
					log.Println("!!!!!! Failed to download Artifact", fileName, " @ ", downloadHref)
				}
			}()
		}
	}
	wg.Wait()

	// Call update procedure to kill player, update music and playlist and startup player again
	if StartUpdate(chunks) != nil {
		log.Println("!! Error updating")
		client.UpdateActionStatus(actionId, hawkbit.STATUS_EXEC_CLOSED, hawkbit.STATUS_RESULT_FAILED)
	} else {
		log.Println("** Done updating.")
		// Tell hawkbit we've completed updating
		client.UpdateActionStatus(actionId, hawkbit.STATUS_EXEC_CLOSED, hawkbit.STATUS_RESULT_SUCCESS)
	}

}

func SaveContent(path string, a *hawkbit.Artifact) error {
	os.MkdirAll("."+string(filepath.Separator)+path, 0777)

	err := ioutil.WriteFile(path+"/"+a.Filename, a.Content, 0644)
	return err
}

type DownloadTask struct {
	Href string

	FileSize int
	Filename string

	ChunkPart    string
	ChunkName    string
	ChunkVersion string
}

func CopyFile(src, dest string, perm os.FileMode) {
	content, err := ioutil.ReadFile(src)
	if err != nil {
		panic(err.Error())
	}
	ioutil.WriteFile(dest, content, perm)
}

func StartUpdate(chunks []*hawkbit.Chunk) error {
	log.Println("@@@@@@ Start Updating..")
	log.Println("@@@@@@ Stop Music Player")
	KillPlayer()

	for _, c := range chunks {
		for _, a := range c.Artifacts {
			src := fmt.Sprintf("./hawkbit/%s/%s/%s/%s", c.Part, c.Name, c.Version, a.Filename)

			if c.Part == "os" {
				dest := "../player/klingenplayer"
				CopyFile(src, dest, 0744)
				log.Println("++++++ Copied Music Player Binary")
			} else if c.Part == "bApp" {
				if strings.HasSuffix(a.Filename, ".m3u") {
					// Copy playlist file
					dest := "../player/playlist"
					CopyFile(src, dest, 0644)

					log.Println("++++++ Copied Playlist File")

				} else if strings.HasSuffix(a.Filename, ".mp3") {
					// Copy music file
					dest := "../player/music/" + a.Filename
					CopyFile(src, dest, 0644)

					log.Println("++++++ Copied Music File ", a.Filename)
				} else {
					log.Println("[!!!!!] Unknown File Type. Only mp3 and m3u files are recognized.")
				}
			} else {
				log.Println("[!!!!!] Unknown Chunk/Part")
			}
		}
	}

	// Delete all downloaded firmware files.
	// TODO: Wait and ensure all copy processes are completed before cleaning up.
	// CleanUp()

	// Start Player again
	log.Println("@@@@@@ Start Music Player")
	StartPlayer()

	return nil
}

func StartPlayer() {
	cmd := exec.Command("../player/klingenplayer")
	file, err := os.Create("out.log")
	if err != nil {
		panic(err.Error())
	}

	cmd.Stdout = file
	err = cmd.Start()

	if err != nil {
		panic(err.Error())
	}

	PLAYER_PROCESS = cmd
	log.Println("@@@@@@ KlingenCloud Player Started")
}

func KillPlayer() {
	if PLAYER_PROCESS != nil {
		if err := PLAYER_PROCESS.Process.Kill(); err != nil {
			log.Fatal("!! Failed to kill Music Player process: ", err)
		} else {
			log.Println("@@ KlingenCloud Player Stopped")
		}
	}
}

func CleanUp() {
	os.RemoveAll("./hawkbit")
}

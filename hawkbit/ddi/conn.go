package ddi

import (
	"github.com/zubairhamed/iot-suite-sdk-go/hawkbit/base"
	. "github.com/zubairhamed/iot-suite-sdk-go/hawkbit"
	"bytes"
	"time"
	"fmt"
	"errors"
	"io"
	"crypto/tls"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"strings"
	"sync"
)

type HawkbitDDIConnection struct {
	base.HawkbitConnection
	endpoint string
	tenant string
	target string
	token string
	configuration *Configuration
	updateStarted bool
}

func (c *HawkbitDDIConnection) WaitForUpdates(ch chan *Message, pollingSecs int) {
	go func() {
		ticker := time.NewTicker(time.Second * time.Duration(pollingSecs))
		for range ticker.C {
			if !c.updateStarted {
				u, a := c.ReadyForUpdate()
				if u {
					msg := &Message{
						ActionId: a,
					}
					ch <- msg
				}
			}
		}
	}()
}

func (c *HawkbitDDIConnection) UpdateActionStatus(actionId string, e ExecStatus, r ResultStatus) (err error) {
	httpUrl := fmt.Sprintf("%s/%s/controller/v1/%s/deploymentBase/%s/feedback", c.endpoint, c.tenant, c.target, actionId)
	var bodyString string
	switch {
	case e == STATUS_EXEC_PROCEEDING:
		c.updateStarted = true
		bodyString = fmt.Sprintf(`{ "id" : "%s", "time" : "%d", "status" : { "result" : { "finished" : "none" }, "execution" : "proceeding" } }`, actionId, time.Now())
		break

	case e == STATUS_EXEC_CLOSED && r == STATUS_RESULT_SUCCESS:
		c.updateStarted = false
		bodyString = fmt.Sprintf(`{ "id" : "%s", "time" : "%d", "status" : { "result" : { "finished" : "success" }, "execution" : "closed" } }`, actionId, time.Now())
		break

	case e == STATUS_EXEC_CLOSED && r == STATUS_RESULT_FAILED:
		c.updateStarted = false
		bodyString = fmt.Sprintf(`{ "id" : "%s", "time" : "%d", "status" : { "result" : { "finished" : "failure" }, "execution" : "closed" } }`, actionId, time.Now())
		break
	}

	req, err := c.createHttpRequest("POST", httpUrl, bytes.NewBufferString(bodyString))
	if err != nil {
		return
	}

	httpClient := c.createHttpClientInstance()
	resp, err := httpClient.Do(req)
	if err != nil {
		return
	}

	if resp.StatusCode != 200 {
		err = errors.New("Status code returned not 200")

		return
	}
	return
}

func (c *HawkbitDDIConnection) GetDownloadableChunks() (chunks []*Chunk, err error){
	d, err := c.GetDeploymentBase()
	if err != nil {
		return
	}

	chunks = d.Deployment.Chunks
	return
}

func (c *HawkbitDDIConnection) GetDeploymentBase() (dbInfo *DeploymentBaseInfo, err error) {
	actions, err := c.GetActions()
	if err != nil {
		return
	}

	if actions.Links.DeploymentBase == nil {
		return
	}

	content, err := c.doGet(actions.Links.DeploymentBase.Href)
	dbInfo, err = JsonToDeploymentBaseResponse(content)
	if err != nil {
		return
	}

	return
}

func (c *HawkbitDDIConnection) DownloadArtifact(a *Artifact) (err error) {
	httpUrl := a.Links.DownloadHttp.Href
	content, err := c.doGet(httpUrl)
	if err != nil {
		return
	}

	a.Content = content

	return
}

func (c *HawkbitDDIConnection) ReadyForUpdate() (updateNow bool, actionId string) {
	updateNow = false
	actions, err := c.GetActions()
	if err != nil {
		updateNow = false

		return
	}

	if actions.Links.DeploymentBase != nil {
		updateNow = true

		href := actions.Links.DeploymentBase.Href
		actionId = href[strings.LastIndex(href, "/")+1 : strings.LastIndex(href, "?")]
	}
	return
}

func (c *HawkbitDDIConnection) DownloadArtifacts(downloadingCh chan Artifact, downloadedCh chan Artifact, parallel bool) (count int, err error){
	chunks, err := c.GetDownloadableChunks()
	if err != nil {
		return
	}
	count = 0
	for _, chunk := range chunks {
		count += len(chunk.Artifacts)
	}

	go func() {
		var wgChunk sync.WaitGroup
		wgChunk.Add(len(chunks))
		for _, chunk := range chunks {
			var wgArtifacts sync.WaitGroup
			wgArtifacts.Add(len(chunk.Artifacts))
			go func() {
				for _, artifact := range chunk.Artifacts {

					// dereference pointer
					artifact := *artifact

					if parallel {
						go func() {
							c.doDownloadArtifact(artifact, downloadingCh, downloadedCh)
							wgArtifacts.Done()
						}()
					} else {
						c.doDownloadArtifact(artifact, downloadingCh, downloadedCh)
						wgArtifacts.Done()
					}
				}
			}()
			wgArtifacts.Wait()
			wgChunk.Done()
		}
		wgChunk.Wait()
		c.updateStarted = false
	}()
	return
}

func (c *HawkbitDDIConnection) doDownloadArtifact(a Artifact, downloadingCh chan Artifact, downloadedCh chan Artifact) {
	downloadingCh <- a
	if c.DownloadArtifact(&a) == nil {
		downloadedCh <- a
	}
}

func (c *HawkbitDDIConnection) GetActions() (ar *Action, err error) {
	req, err := c.createHttpRequest("GET", fmt.Sprintf("%s/%s/controller/v1/%s", c.endpoint, c.tenant, c.target), nil)
	if err != nil {
		return nil, err
	}

	httpClient := c.createHttpClientInstance()
	resp, err := httpClient.Do(req)
	if err != nil {
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &ar)

	return
}

func (c *HawkbitDDIConnection) createHttpRequest(method, url string, body io.Reader) (req *http.Request, err error) {
	req, err = http.NewRequest(method, url, body)
	if err == nil {
		req.Header.Set("Authorization", "TargetToken "+c.token)
		req.Header.Set("Content-Type", "application/json")
	}
	return
}

func (c *HawkbitDDIConnection) createHttpClientInstance() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		//Proxy:           http.ProxyURL(c.proxy),
	}

	return &http.Client{Transport: tr}
}

func (c *HawkbitDDIConnection) doGet(s string) (content []byte, err error) {
	req, err := c.createHttpRequest("GET", s, nil)
	if err != nil {
		return
	}

	httpClient := c.createHttpClientInstance()
	resp, err := httpClient.Do(req)
	if err != nil {
		return
	}

	content, err = ioutil.ReadAll(resp.Body)

	return
}




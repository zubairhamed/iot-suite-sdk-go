package ddi

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
	. "github.com/zubairhamed/iot-suite-sdk-go/hawkbit"
)

func NewDefaultHawkbitClient(server, tenant, target, targetToken string) HawkbitClient {
	return &HawkbitDDIClient{
		server:        server,
		tenant:        tenant,
		target:        target,
		targetToken:   targetToken,
		updateStarted: false,
	}
}

type HawkbitDDIClient struct {
	proxy       *url.URL
	server      string
	tenant      string
	target      string
	targetToken string

	updateStarted bool

	fnUpdate CallbackFn
}

func (c *HawkbitDDIClient) UseProxy(s string) error {
	p, err := url.Parse(s)
	c.proxy = p

	return err
}

func (c *HawkbitDDIClient) Get(s string) (content []byte, err error) {
	req, err := c.createHttpRequest("GET", s, nil)
	if err != nil {
		return
	}

	httpClient := c.createHttpClientInstance()
	resp, err := httpClient.Do(req)

	if err != nil {
		log.Println(err)
		return
	}

	content, err = ioutil.ReadAll(resp.Body)

	return
}

func (c *HawkbitDDIClient) GetActions() (ar *Action, err error) {
	req, err := c.createHttpRequest("GET", fmt.Sprintf("%s/%s/controller/v1/%s", c.server, c.tenant, c.target), nil)
	if err != nil {
		return nil, err
	}

	httpClient := c.createHttpClientInstance()
	resp, err := httpClient.Do(req)
	if err != nil {
		log.Println(err)
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &ar)

	return
}

func (c *HawkbitDDIClient) createHttpRequest(method, url string, body io.Reader) (req *http.Request, err error) {
	req, err = http.NewRequest(method, url, body)
	if err == nil {
		req.Header.Set("Authorization", "TargetToken "+c.targetToken)
		req.Header.Set("Content-Type", "application/json")
	}
	return
}

func (c *HawkbitDDIClient) createHttpClientInstance() *http.Client {
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
		Proxy:           http.ProxyURL(c.proxy),
	}

	return &http.Client{Transport: tr}
}

func (c *HawkbitDDIClient) ShouldUpdateNow() (updateNow bool, actionId string) {
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

func (c *HawkbitDDIClient) GetDownloadableChunks() (chunks []*Chunk) {
	d, err := c.GetDeploymentBase()
	if err != nil {
		log.Println(err)
	}
	return d.Deployment.Chunks
}

func (c *HawkbitDDIClient) GetDeploymentBase() (dbInfo *DeploymentBaseInfo, err error) {
	actions, err := c.GetActions()
	if err != nil {
		return
	}

	if actions.Links.DeploymentBase == nil {
		return
	}

	content, err := c.Get(actions.Links.DeploymentBase.Href)
	dbInfo, err = JsonToDeploymentBaseResponse(content)
	if err != nil {
		return
	}

	return
}

func (c *HawkbitDDIClient) DownloadArtifact(a *Artifact) (err error) {
	httpUrl := a.Links.DownloadHttp.Href

	content, err := c.Get(httpUrl)
	if err != nil {
		return
	}

	a.Content = content

	return
}

func (c *HawkbitDDIClient) Start() {
	ticker := time.NewTicker(time.Millisecond * 3000)
	for range ticker.C {
		if !c.updateStarted {
			u, a := c.ShouldUpdateNow()
			if u {
				c.fnUpdate(c, a)
			}
		}
	}
}

func (c *HawkbitDDIClient) OnUpdate(fn CallbackFn) {
	c.fnUpdate = fn
}

func (c *HawkbitDDIClient) UpdateActionStatus(actionId string, e ExecStatus, r ResultStatus) (err error) {

	httpUrl := fmt.Sprintf("%s/%s/controller/v1/%s/deploymentBase/%s/feedback", c.server, c.tenant, c.target, actionId)
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

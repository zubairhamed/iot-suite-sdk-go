package hawkbit

import "strings"

type Action struct {
	Config struct {
		Polling struct {
			Sleep string `json:"sleep"`
		} `json:"polling"`
	} `json:"config"`

	Links struct {
		DeploymentBase *Href `json:"deploymentBase"`
		ConfigData     *Href `json:"configData"`
	} `json:"_links"`
}

type Href struct {
	Href string `json:"href"`
}

func (h *Href) GetActionId() string {
	href := h.Href

	return href[strings.LastIndex(href, "/"):]
}

type DeploymentBaseInfo struct {
	Deployment struct {
		Download string   `json:"download"`
		Update   string   `json:"update"`
		Chunks   []*Chunk `json:"chunks"`
	} `json:"deployment"`
}

type Chunk struct {
	Part      string      `json:"part"`
	Version   string      `json:"version"`
	Name      string      `json:"name"`
	Artifacts []*Artifact `json:"artifacts"`
}

type Artifact struct {
	Filename string `json:"filename"`
	Hashes   struct {
		SHA1 string `json:"sha1"`
		MD5  string `json:"md5"`
	} `json: "hashes"`
	Size  int `json:"size"`
	Links struct {
		DownloadHttp *Href `json:"download-http"`
		Md5SumHttp   *Href `json:"md5sum-http"`
		Download     *Href `json:"download"`
		Md5Sum       *Href `json:"md5sum"`
	} `json:"_links"`
	Content []byte
}

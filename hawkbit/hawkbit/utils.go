package hawkbit

import "encoding/json"

func JsonToDeploymentBaseResponse(content []byte) (resp *DeploymentBaseInfo, err error) {
	err = json.Unmarshal(content, &resp)

	return
}

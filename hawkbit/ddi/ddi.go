package ddi

import (
	. "github.com/zubairhamed/iot-suite-sdk-go/hawkbit"
)

func Dial(endpoint, tenant, target, token string, cfg *Configuration) (*HawkbitDDIConnection, error) {
	return &HawkbitDDIConnection{
		endpoint: endpoint,
		tenant: tenant,
		target: target,
		token: token,
		configuration: cfg,
	}, nil
}

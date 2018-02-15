package examples

import (
	"github.com/zubairhamed/iot-suite-sdk-go/things/client"
)

// For IoT Things
var ENDPOINT_URL_REST = "https://things.apps.bosch-iot-cloud.com"
var ENDPOINT_URL_WS = "wss://things.apps.bosch-iot-cloud.com"
var USERNAME = ""
var PASSWORD = ""
var APITOKEN = ""
var DEFAULT_CLIENT_CONFIG = &client.Configuration{
	SkipSslVerify: true,
}

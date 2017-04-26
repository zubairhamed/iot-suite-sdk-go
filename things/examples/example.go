package examples

import (
	"github.com/zubairhamed/iot-suite-sdk-go/things/client"
)

// For IoT Things
var ENDPOINT_URL_REST = "https://things.apps.bosch-iot-cloud.com"
var ENDPOINT_URL_WS = "wss://things.apps.bosch-iot-cloud.com"
var USERNAME = "GoTest"
var PASSWORD = "GoTestPw1!"
var APITOKEN = "b7dca5d42d1245d6b5ee06e3b3f1a044"
var DEFAULT_CLIENT_CONFIG = &client.Configuration{
	SkipSslVerify: true,
}

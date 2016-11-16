package examples

import "github.com/zubairhamed/iot-suite-sdk-go/things/client"

// For IoT Things
var ENDPOINT_URL_REST = "https://things-int.apps.bosch-iot-cloud.com"
var ENDPOINT_URL_WS = "wss://things-int.apps.bosch-iot-cloud.com"
var USERNAME = "Zubair"
var PASSWORD = "ZubairPw1!"
var APITOKEN = "c8746ff31faf46dabf68eb7188df1694"
var DEFAULT_CLIENT_CONFIG = &client.Configuration{
	SkipSslVerify: true,
}
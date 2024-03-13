package model

const (
	HeaderKeyUserAgent        = "User-Agent"
	HeaderKeyUserIP           = "X-Forwarded-For"
	HeaderKeyRequestID        = "requestid"
	HeaderKeyDeviceModel      = "dp-device-model"
	HeaderKeyDeviceOs         = "dp-device-os"
	HeaderKeyDeviceOsVersion  = "DP-Device-Os-Version"
	HeaderKeyDeviceAppVersion = "Dp-App-Version"
)

const (
	LoggerKeyRequestID        = "REQUEST_ID"
	LoggerKeyOperation        = "OPERATION"
	LoggerKeyUserIP           = "USER_IP"
	LoggerKeyUserAgent        = "USER_AGENT"
	LoggerKeyDeviceModel      = "DEVICE_MODEL"
	LoggerKeyDeviceOs         = "DEVICE_OS"
	LoggerKeyDeviceOsVersion  = "DEVICE_OS_VERSION"
	LoggerKeyDeviceAppVersion = "APP_VERSION"
	ContextLogger             = "contextLogger"
	ContextHeader             = "contextHeader"
)

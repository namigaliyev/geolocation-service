package middleware

import (
	"context"
	"geolocation-service/internal/model"
	log "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

var headers = []string{
	"x-request-id",
	"x-b3-traceid",
	"x-b3-spanid",
	"x-b3-parentspanid",
	"x-b3-sampled",
	"x-b3-flags",
	"x-ot-span-context",
	"User-Agent",
	"X-Forwarded-For",
	"requestid",
}

func RequestParamsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		ctx := r.Context()

		requestID := r.Header.Get(model.HeaderKeyRequestID)
		operation := r.RequestURI
		userAgent := r.Header.Get(model.HeaderKeyUserAgent)
		userIP := r.Header.Get(model.HeaderKeyUserIP)
		deviceModel := r.Header.Get(model.HeaderKeyDeviceModel)
		deviceOs := r.Header.Get(model.HeaderKeyDeviceOs)
		deviceOsVersion := r.Header.Get(model.HeaderKeyDeviceOsVersion)
		appVersion := r.Header.Get(model.HeaderKeyDeviceAppVersion)

		if len(requestID) == 0 {
			requestID = primitive.NewObjectID().Hex()
		}
		fields := log.Fields{}
		addLoggerParam(fields, model.LoggerKeyRequestID, requestID)
		addLoggerParam(fields, model.LoggerKeyOperation, operation)
		addLoggerParam(fields, model.LoggerKeyUserAgent, userAgent)
		addLoggerParam(fields, model.LoggerKeyUserIP, userIP)
		addLoggerParam(fields, model.LoggerKeyDeviceModel, deviceModel)
		addLoggerParam(fields, model.LoggerKeyDeviceOs, deviceOs)
		addLoggerParam(fields, model.LoggerKeyDeviceOsVersion, deviceOsVersion)
		addLoggerParam(fields, model.LoggerKeyDeviceAppVersion, appVersion)

		logger := log.WithFields(fields)
		header := http.Header{}

		for _, v := range headers {
			header.Add(v, r.Header.Get(v))
		}

		ctx = context.WithValue(ctx, model.ContextLogger, logger)
		ctx = context.WithValue(ctx, model.ContextHeader, header)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func addLoggerParam(fields log.Fields, field string, value string) {
	if len(value) > 0 {
		fields[field] = value
	}
}

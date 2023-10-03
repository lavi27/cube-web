package middlewares

import (
	"cubeWeb/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
)

func Logger() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		startTime := time.Now()

		ctx.Next()

		endTime := time.Now()
		latencyTime := endTime.Sub(startTime)
		reqMethod := ctx.Request.Method
		reqUri := ctx.Request.RequestURI
		statusCode := ctx.Writer.Status()
		clientIP := ctx.ClientIP()

		logger := utils.HTTPLogger.WithFields(logrus.Fields{
			"METHOD":    reqMethod,
			"URI":       reqUri,
			"STATUS":    statusCode,
			"LATENCY":   latencyTime,
			"CLIENT_IP": clientIP,
		})

		if statusCode == 200 || statusCode == 400 {
			logger.Info("HTTP_REQ")
		} else if statusCode == 500 {
			logger.Error("HTTP_REQ")
		} else {
			logger.Warn("HTTP_REQ")
		}

		ctx.Next()
	}
}

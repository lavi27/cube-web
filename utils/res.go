package utils

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func ResOK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, data)
}

func ResError(c *gin.Context, status int, errCode int, msg string) {
	c.AbortWithStatusJSON(status, gin.H{
		"status":    status,
		"errorCode": errCode,
		"error":     HttpStatusName[status],
		"message":   msg,
	})
}

func InternalError(c *gin.Context, err error) {
	ResError(c, http.StatusInternalServerError, 200, "Unknown server error occured")
	ErrLogger.Error(err)
}

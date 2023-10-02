package utils

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func ResOK(c *gin.Context, data any) {
	c.JSON(http.StatusOK, data)
}

func ResError(c *gin.Context, status int, errCode int, msg string) {
	c.JSON(status, gin.H{
		"status":    status,
		"errorCode": errCode,
		"error":     HttpStatusName[status],
		"message":   msg,
	})
}

func InternalError(c *gin.Context, err error) {
	ResError(c, http.StatusInternalServerError, 1, "Unknown server error occured")
	log.Fatal(err)
}

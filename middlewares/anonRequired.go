package middlewares

import (
	"cubeWeb/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AnonRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if utils.IsLoggedIn(c) {
			utils.ResError(c, http.StatusBadRequest, 101, "You already signed in")
			return
		}

		c.Next()
	}
}

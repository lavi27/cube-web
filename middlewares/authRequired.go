package middlewares

import (
	"cubeWeb/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !utils.IsLoggedIn(c) {
			utils.ResError(c, http.StatusBadRequest, 100, "You must be signed in")
			return
		}

		c.Next()
	}
}

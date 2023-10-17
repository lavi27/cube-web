package middlewares

import (
	"cubeWeb/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckSessionId() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := utils.GetUserId(c)
		if err != nil {
			utils.ResError(c, http.StatusGone, 102, "Invaild sessionId")
			return
		}

		c.Next()
	}
}

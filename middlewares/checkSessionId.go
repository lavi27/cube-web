package middlewares

import (
	"cubeWeb/env"
	"cubeWeb/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckSessionId() gin.HandlerFunc {
	return func(c *gin.Context) {
		_, err := c.Cookie("sessionId")
		if err != nil {
			utils.ResError(c, http.StatusRequestTimeout, 102, "Cookie 'sessionId' not found.")
			return
		}

		_, err = utils.GetUserId(c)
		if err != nil {
			c.SetCookie("sessionId", "", -1, "/", env.ClientIP, false, true)
			utils.ResError(c, http.StatusRequestTimeout, 102, err.Error())
			return
		}

		c.Next()
	}
}

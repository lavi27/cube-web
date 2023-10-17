package middlewares

import (
	"cubeWeb/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckReqBody[T comparable]() gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody T
		if err := c.BindJSON(&reqBody); err != nil {
			utils.ResError(c, http.StatusBadRequest, 103, "Invaild request body")
			return
		}

		c.Next()
	}
}

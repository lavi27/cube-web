package middlewares

import (
	"bytes"
	"cubeWeb/utils"
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

func CheckReqBody[T comparable]() gin.HandlerFunc {
	return func(c *gin.Context) {
		var reqBody T

		bodyBytes, err := c.GetRawData()
		if err != nil {
			utils.ResError(c, http.StatusBadRequest, 103, "Invaild request body")
			return
		}
		c.Request.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))

		if err := json.Unmarshal(bodyBytes, &reqBody); err != nil {
			utils.ResError(c, http.StatusBadRequest, 103, "Invaild request body")
			return
		}

		c.Next()
	}
}

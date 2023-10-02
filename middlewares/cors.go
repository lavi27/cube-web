package middlewares

import "github.com/gin-gonic/gin"

func Cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		// method := c.Request.Method
		origin := c.GetHeader("Origin")

		c.Header("Access-Control-Allow-Origin", origin)
		c.Header("Access-Control-Allow-Credentials", "true")
		c.Header("Access-Control-Allow-Headers", "Access-Control-Allow-Headers,Cookie, Origin, X-Requested-With, Content-Type, Accept, Authorization, Token, Timestamp, UserId")
		c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS,DELETE,PUT")
		c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Content-Type,cache-control")

		// if method == "OPTIONS" {
		// 	//c.AbortWithStatus(http.StatusNoContent)
		// 	c.AbortWithStatus(http.StatusOK)
		// }

		c.Next()
	}
}
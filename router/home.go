// package main

// import (
//   "net/http"
//   "fmt"
//   "os"
//   "log"
  
//   "github.com/joho/godotenv"
//   "github.com/gin-gonic/gin"
// )

// func init() {
//   err := godotenv.Load(".env")
    
//   if err != nil {
//       log.Fatal("Error loading .env file")
//   }
  
//   // os.Getenv("DBNAME")
// }

// func SetupRouter() *gin.Engine {
// 	r := gin.Default()

// 	r.GET("/", func(c *gin.Context) {
// 		c.String(http.StatusOK, "pong")

//     c.JSON()
// 	})

// 	r.GET("/:name", func(c *gin.Context) {
// 		user := c.Params.ByName("name")
// 		value, ok := db[user]
// 		if ok {
// 			c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
// 		} else {
// 			c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
// 		}
// 	})

// 	return r
// }
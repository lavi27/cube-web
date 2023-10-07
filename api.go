package main

import (
	"cubeWeb/env"
	"cubeWeb/middlewares"
	"cubeWeb/model"
	"cubeWeb/router"

	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.ForwardedByClientIP = true
	r.SetTrustedProxies([]string{env.ClientIP})
	r.Static("/static", "./static")

	r.Use(gin.Recovery())
	r.Use(middlewares.Logger())
	r.Use(middlewares.Cors())
	r.Use(middlewares.Session())

	api := r.Group("/api")
	router.SetPostRouter(api)
	router.SetUserRouter(api)

	return r
}

func main() {
	model.ConnectDB()
	r := setupRouter()

	r.Run(":8080")
}

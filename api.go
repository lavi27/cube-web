package main

import (
	"cubeWeb/env"
	"cubeWeb/middlewares"
	"cubeWeb/model"
	"cubeWeb/router"
	"cubeWeb/utils"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

func setupRouter() *gin.Engine {
	if utils.IsDebug {
		gin.SetMode(gin.DebugMode)
	} else {
		gin.SetMode(gin.ReleaseMode)
	}

	r := gin.Default()

	r.ForwardedByClientIP = true
	r.SetTrustedProxies([]string{env.ClientIP})
	r.Static("/static", "./static")

	r.Use(gin.Recovery())
	r.Use(middlewares.Logger())
	r.Use(middlewares.Cors())
	r.Use(sessions.Sessions("session", model.SessionStore))

	api := r.Group("/api")
	router.SetPostRouter(api)
	router.SetUserRouter(api)
	router.SetAccountRouter(api)

	return r
}

func main() {
	model.ConnectDB()
	r := setupRouter()

	r.Run(":8080")
}

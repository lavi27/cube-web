package utils

import (
	"cubeWeb/env"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const IsDebug = true
const SessionTimeoutSec = 60 * 60 * 3 //3h

var IsProduction = env.Enviroment == "production" || env.Enviroment != "devlopment"
var IsDevlopment = env.Enviroment == "devlopment"

func GetSession(c *gin.Context) sessions.Session {
	session := sessions.Default(c)
	session.Options(sessions.Options{MaxAge: SessionTimeoutSec})
	return session
}

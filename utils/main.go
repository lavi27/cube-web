package utils

import (
	"cubeWeb/env"
	"reflect"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

const IsDebug = false
const SessionTimeoutSec = 60 * 60 * 3 //3h

var IsProduction = env.Enviroment == "production" || env.Enviroment != "devlopment"
var IsDevlopment = env.Enviroment == "devlopment"

func GetSession(c *gin.Context) sessions.Session {
	session := sessions.Default(c)
	session.Options(sessions.Options{MaxAge: SessionTimeoutSec})
	return session
}

func GetUserIdBySessionId(session sessions.Session, sessionId string) int64 {
	realSession := reflect.ValueOf(session).Elem().FieldByName("session")

	if !realSession.IsValid() {
		panic("")
	}

	realSession = realSession.Elem().FieldByName(sessionId)

	if !realSession.IsValid() {
		panic("")
	}

	return realSession.Int()
}

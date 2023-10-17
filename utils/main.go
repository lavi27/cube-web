package utils

import (
	"cubeWeb/env"
	"errors"
	"reflect"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var IsProduction = env.Enviroment == "production" || env.Enviroment != "devlopment"
var IsDevlopment = env.Enviroment == "devlopment"

func GetSession(c *gin.Context) sessions.Session {
	session := sessions.Default(c)
	session.Options(sessions.Options{MaxAge: SessionTimeoutSec})
	return session
}

func IsLoggedIn(c *gin.Context) bool {
	_, err := c.Cookie("sessionId")

	if err != nil {
		return true
	} else {
		return false
	}
}

func GetUserId(c *gin.Context) (int, error) {
	session := GetSession(c)
	sessionId, _ := c.Cookie("sessionId")

	realSession := reflect.ValueOf(session).Elem().FieldByName("session")

	if !realSession.IsValid() {
		return 0, errors.New("session not found")
	}

	realSession = realSession.Elem().FieldByName(sessionId)

	if !realSession.IsValid() {
		return 0, errors.New("sessionId not found")
	}

	return int(realSession.Int()), nil
}

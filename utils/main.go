package utils

import (
	"cubeWeb/env"
	"encoding/json"
	"errors"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var IsProduction = env.Enviroment == "production" || env.Enviroment != "devlopment"
var IsDevlopment = env.Enviroment == "devlopment"

func GetSession(c *gin.Context) sessions.Session {
	session := sessions.Default(c)
	session.Options(sessions.Options{MaxAge: SessionTimeoutSec, Path: "/"})
	return session
}

func IsLoggedIn(c *gin.Context) bool {
	_, err := c.Cookie("sessionId")

	if err != nil {
		return false
	} else {
		return true
	}
}

func GetUserId(c *gin.Context) (int, error) {
	session := GetSession(c)
	sessionId, err := c.Cookie("sessionId")
	if err != nil {
		return 0, ErrNotFound
	}

	val := session.Get(sessionId)
	if val == nil {
		return 0, errors.New("invaild sessionId")
	}

	valInt, ok := val.(int)
	if !ok {
		return 0, errors.New("invaild sessionId value type")
	}

	return valInt, nil
}

func GetBodyJSON(c *gin.Context, ptr any) {
	bodyBytes, _ := c.GetRawData()
	_ = json.Unmarshal(bodyBytes, ptr)
}

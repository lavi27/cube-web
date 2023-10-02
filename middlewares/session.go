package middlewares

import (
	"cubeWeb/env"
	"database/sql"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/postgres"
	"github.com/gin-gonic/gin"
)

func Session() gin.HandlerFunc {
	dsn := "host=" + env.DBIP + " user=" + env.DBUser + " password=" + env.DBPassword + " dbname=" + env.DBName + " port=" + env.DBPort + " sslmode=disable"
	db, err := sql.Open("postgres", dsn)
	if err != nil {
		panic("세션 db 연결 실패")
	}

	store, err := postgres.NewStore(db)
	if err != nil {
		panic(err.Error())
	}

	return sessions.Sessions("session", store)
}

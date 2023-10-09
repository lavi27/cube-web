package router

import (
	"cubeWeb/model"
	"cubeWeb/utils"
	"net/http"

	"github.com/gin-gonic/gin"
)

func SetAccountRouter(rg *gin.RouterGroup) {
	group := rg.Group("/account")
	group.GET("/", getAccount)
}

func getAccount(c *gin.Context) {
	session := utils.GetSession(c)

	sessionId, err := c.Cookie("sessionId")
	if err != nil {
		utils.ResError(c, http.StatusBadRequest, 1, "You must be logged in")
		return
	}

	userId := utils.GetUserIdBySessionId(session, sessionId)

	var dbQuery model.User
	if err := model.DB.Where(
		model.User{UserId: int(userId)},
	).Find(&dbQuery).Error; err != nil {
		utils.InternalError(c, err)
		return
	}

	if dbQuery.UserId == 0 {
		utils.ResError(c, http.StatusNotFound, 2, "Account not found")
		return
	}

	utils.ResOK(c, int(userId))
}

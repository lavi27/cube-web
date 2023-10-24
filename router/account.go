package router

import (
	"cubeWeb/env"
	"cubeWeb/middlewares"
	"cubeWeb/model"
	"cubeWeb/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func SetAccountRouter(rg *gin.RouterGroup) {
	group := rg.Group("/account")

	group.GET("/",
		middlewares.CheckSessionId(),
		getAccount,
	)
	group.POST("/signin/",
		middlewares.AnonRequired(),
		middlewares.CheckReqBody[signinReqBody](),
		postSignin,
	)
	group.POST("/signup/",
		middlewares.AnonRequired(),
		middlewares.CheckReqBody[signupReqBody](),
		postSignup,
	)
}

func getAccount(c *gin.Context) {
	userId, _ := utils.GetUserId(c)

	var dbQuery model.User
	if err := model.DB.Where(
		model.User{UserId: userId},
	).Find(&dbQuery).Error; err != nil {
		utils.InternalError(c, err)
		return
	}

	if dbQuery.UserId == 0 {
		c.SetCookie("sessionId", "", -1, "/", env.ClientIP, false, true)
		utils.ResError(c, http.StatusRequestTimeout, 102, "Invaild sessionId")
		return
	}

	utils.ResOK(c, userId)
}

type signinReqBody struct {
	UserName string `json:"userName"`
	UserPW   string `json:"userPassword"`
}

func postSignin(c *gin.Context) {
	session := utils.GetSession(c)
	var reqBody signinReqBody
	utils.GetBodyJSON(c, &reqBody)

	var dbQuery model.User
	if err := model.DB.Where(
		model.User{UserNm: reqBody.UserName},
	).Find(&dbQuery).Error; err != nil {
		utils.InternalError(c, err)
		return
	}

	if dbQuery.UserId == 0 {
		utils.ResError(c, http.StatusBadRequest, 3, "Failed to signin")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbQuery.UserPw), []byte(reqBody.UserPW)); err != nil {
		utils.ResError(c, http.StatusBadRequest, 3, "Failed to signin")
		return
	}

	uuid := uuid.New().String()
	session.Set(uuid, dbQuery.UserId)
	session.Save()
	c.SetCookie("sessionId", uuid, utils.SessionTimeoutSec, "/", env.ClientIP, false, true)

	utils.ResOK(c, nil)
}

type signupReqBody struct {
	UserName string `json:"userName"`
	UserPW   string `json:"userPassword"`
}

func postSignup(c *gin.Context) {
	session := utils.GetSession(c)
	var reqBody signupReqBody
	utils.GetBodyJSON(c, &reqBody)

	var dbQuery model.User
	if err := model.DB.Where(
		model.User{UserNm: reqBody.UserName},
	).Find(&dbQuery).Error; err != nil {
		utils.InternalError(c, err)
		return
	}

	if dbQuery.UserId != 0 {
		utils.ResError(c, http.StatusBadRequest, 3, "Same userName already exists")
		return
	}

	pwHash, err := bcrypt.GenerateFromPassword([]byte(reqBody.UserPW), bcrypt.MinCost)
	if err != nil {
		utils.InternalError(c, err)
		return
	}

	user := model.User{
		UserNm: reqBody.UserName,
		UserPw: string(pwHash),
	}
	if err := model.DB.Create(&user).Error; err != nil {
		utils.InternalError(c, err)
		return
	}

	uuid := uuid.New().String()
	session.Set(uuid, user.UserId)
	session.Save()
	c.SetCookie("sessionId", uuid, utils.SessionTimeoutSec, "/", env.ClientIP, false, true)

	utils.ResOK(c, nil)
}

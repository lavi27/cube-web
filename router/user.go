package router

import (
	"cubeWeb/env"
	"cubeWeb/model"
	"cubeWeb/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

func SetUserRouter(rg *gin.RouterGroup) {
	ping := rg.Group("/user")
	ping.GET("/", getUser)
	ping.POST("/signin", postSignin)
	ping.POST("/signup", postSignup)
}

type userResData struct {
	UserId         int    `json:"userId,omitempty"`
	UserName       string `json:"userName,omitempty"`
	UserNickname   string `json:"userNickname,omitempty"`
	CreateDate     int64  `json:"createDate,omitempty"`
	FollowerCount  int    `json:"followerCount"`
	FollowingCount int    `json:"followingCount"`
	PostCount      int    `json:"postCount"`
}

func getUser(c *gin.Context) {
	userIdQur := c.Query("userId")

	userId, err := strconv.Atoi(userIdQur)
	if err != nil {
		utils.ResError(c, http.StatusBadRequest, 1, "Invaild userId query")
		return
	}

	var dbQuery model.User
	if err = model.DB.Model(&model.User{}).Select(
		"users.user_id", "users.user_nm", "users.create_dt", "users.follower_cnt", "users.following_cnt", "users.user_nick_nm",
		"count(posts) as post_cnt",
	).Joins(
		"left join posts on users.user_id = posts.user_id",
	).Group("users.user_id").Where(
		"users.user_id = ?", userId,
	).Find(&dbQuery).Error; err != nil {
		utils.InternalError(c, err)
		return
	}

	if dbQuery.UserId == 0 {
		utils.ResError(c, http.StatusNotFound, 2, "User not found")
		return
	}

	utils.ResOK(c,
		userResData{
			dbQuery.UserId,
			dbQuery.UserNm,
			dbQuery.UserNickNm,
			dbQuery.CreateDt.Unix(),
			dbQuery.FollowerCnt,
			dbQuery.FollowingCnt,
			dbQuery.PostCnt,
		})
}

type signinReqBody struct {
	userName     string
	userPassword string
}

func postSignin(c *gin.Context) {
	session := utils.GetSession(c)
	var reqBody signinReqBody

	if _, err := c.Cookie("sessionId"); err == nil {
		utils.ResError(c, http.StatusBadRequest, 1, "You already logged in")
		return
	}

	if err := c.BindJSON(&reqBody); err != nil {
		utils.ResError(c, http.StatusBadRequest, 2, "Invaild request body")
		return
	}

	var dbQuery model.User
	if err := model.DB.Where(
		model.User{UserNm: reqBody.userName},
	).Find(&dbQuery).Error; err != nil {
		utils.InternalError(c, err)
		return
	}

	if dbQuery.UserId == 0 {
		utils.ResError(c, http.StatusBadRequest, 3, "Failed to signin")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(dbQuery.UserPw), []byte(reqBody.userPassword)); err != nil {
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
	userName     string
	userPassword string
}

func postSignup(c *gin.Context) {
	session := utils.GetSession(c)
	var reqBody signupReqBody

	if _, err := c.Cookie("sessionId"); err == nil {
		utils.ResError(c, http.StatusBadRequest, 1, "You already logged in")
		return
	}

	if err := c.BindJSON(&reqBody); err != nil {
		utils.ResError(c, http.StatusBadRequest, 2, "Invaild request body")
		return
	}

	pwHash, err := bcrypt.GenerateFromPassword([]byte(reqBody.userPassword), bcrypt.MinCost)
	if err != nil {
		utils.ResError(c, http.StatusInternalServerError, 3, "Unknown server error occured")
		return
	}

	user := model.User{
		UserNm: reqBody.userName,
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

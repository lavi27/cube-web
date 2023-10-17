package router

import (
	"cubeWeb/middlewares"
	"cubeWeb/model"
	"cubeWeb/utils"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func SetUserRouter(rg *gin.RouterGroup) {
	group := rg.Group("/user")
	group.GET("/", getUser)
	group.POST("/follow/",
		middlewares.AuthRequired(),
		middlewares.CheckReqBody[followReqBody](),
		middlewares.CheckSessionId(),
		postFollow,
	)
	group.POST("/unfollow/",
		middlewares.AuthRequired(),
		middlewares.CheckReqBody[followReqBody](),
		middlewares.CheckSessionId(),
		postUnfollow,
	)
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
		"users.user_id", "users.user_nm", "users.create_dt",
		"users.follower_cnt", "users.following_cnt", "users.user_nick_nm",
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

type followReqBody struct {
	UserId int `json:"userId"`
}

func postFollow(c *gin.Context) {
	var reqBody followReqBody
	c.BindJSON(&reqBody)

	userId, _ := utils.GetUserId(c)

	var isExist bool

	err := model.DB.Model(
		&model.UserFollow{},
	).Select("count(*) > 0").Where(
		"user_id = ? AND target_user_id = ?",
		userId,
		reqBody.UserId,
	).Find(&isExist).Error
	if err != nil {
		utils.InternalError(c, err)
		return
	}

	if isExist {
		utils.ResError(c, http.StatusBadRequest, 4, "You already followed")
	}

	if err := model.DB.Create(
		model.UserFollow{
			UserId:       userId,
			TargetUserId: reqBody.UserId,
		},
	).Error; err != nil {
		utils.InternalError(c, err)
		return
	}

	utils.ResOK(c, nil)
}

func postUnfollow(c *gin.Context) {
	var reqBody followReqBody
	c.BindJSON(&reqBody)

	userId, _ := utils.GetUserId(c)

	var isExist bool

	err := model.DB.Model(
		&model.UserFollow{},
	).Select("count(*) > 0").Where(
		"user_id = ? AND target_user_id = ?",
		userId,
		reqBody.UserId,
	).Find(&isExist).Error
	if err != nil {
		utils.InternalError(c, err)
		return
	}

	if !isExist {
		utils.ResError(c, http.StatusBadRequest, 5, "You did not followed")
	}

	if err := model.DB.Delete(
		model.UserFollow{
			UserId:       userId,
			TargetUserId: reqBody.UserId,
		},
	).Error; err != nil {
		utils.InternalError(c, err)
		return
	}

	utils.ResOK(c, nil)
}

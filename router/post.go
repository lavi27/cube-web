package router

import (
	"cubeWeb/model"
	"cubeWeb/utils"
	"net/http"
	"reflect"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
)

func SetPostRouter(rg *gin.RouterGroup) {
	ping := rg.Group("/post")
	ping.GET("/", getPost)
	ping.POST("/", postPost)
	ping.GET("/search", getSearch)
}

func getPost(c *gin.Context) {
	postIdQur := c.Query("postId")

	//NOTE - Query PostId
	postId, err := strconv.Atoi(postIdQur)
	if err != nil {
		utils.ResError(c, http.StatusBadRequest, 1, "Invaild postId query")
		return
	}

	post := model.Post{PostId: -1}
	if err = model.DB.Model(&model.Post{}).Select(
		"posts.*",
		"users.user_nm", "users.user_nick_nm",
	).Joins(
		"left join users on posts.user_id = users.user_id",
	).Group("posts.post_id").Find(&post, postId).Error; err != nil {
		utils.InternalError(c, err)
	}

	if post.PostId == -1 {
		utils.ResError(c, http.StatusBadRequest, 2, "Post not found")
		return
	}

	utils.ResOK(c, post)
}

type postReqBody struct {
	sessionId string
	content   string
}

func postPost(c *gin.Context) {
	session := utils.GetSession(c)
	var reqBody postReqBody

	if err := c.BindJSON(&reqBody); err != nil {
		utils.ResError(c, http.StatusBadRequest, 1, "Invaild request body")
	}

	sessionId, err := c.Cookie("sessionId")
	if err != nil {
		utils.ResError(c, http.StatusBadRequest, 2, "You must be logged in")
		return
	}

	realSession := reflect.ValueOf(session).Elem().FieldByName("session")
	userId := realSession.Elem().FieldByName(sessionId).Int()

	post := model.Post{
		UserId:  int(userId),
		Content: reqBody.content,
	}
	if err := model.DB.Create(&post).Error; err != nil {
		utils.InternalError(c, err)
	}

	utils.ResOK(c, nil)
}

func getSearch(c *gin.Context) {
	userId := c.DefaultQuery("userId", "-1")
	dateFromQur := c.DefaultQuery("dateFrom", "-1")
	lengthQur := c.DefaultQuery("length", "10")

	//NOTE - Query UserId
	if _, err := strconv.Atoi(userId); err != nil {
		utils.ResError(c, http.StatusBadRequest, 1, "Invaild userId query")
		return
	}

	//NOTE - Query DateFrom
	var dateFrom time.Time
	if dateFromQur == "-1" {
		dateFrom = time.Now()
	} else {
		i, err := strconv.ParseInt(dateFromQur, 10, 64)
		if err != nil {
			utils.ResError(c, http.StatusBadRequest, 2, "Invaild dateFrom query")
			return
		}
		dateFrom = time.Unix(i, 0)
	}

	//NOTE - Query Length
	length, err := strconv.Atoi(lengthQur)
	if err != nil {
		utils.ResError(c, http.StatusBadRequest, 3, "Invaild length query")
		return
	}
	if length > 100 {
		utils.ResError(c, http.StatusBadRequest, 4, "Too much length, should be under 100")
		return
	}

	//NOTE - DB Select
	var posts []model.Post

	dbQuery := model.DB
	if userId != "-1" {
		dbQuery.Where("user_id = ?", userId)
	}
	err = dbQuery.Where("create_dt < ?", dateFrom).Order("create_dt").Limit(length).Find(&posts).Error

	if err != nil {
		utils.InternalError(c, err)
	}

	utils.ResOK(c, posts)
}

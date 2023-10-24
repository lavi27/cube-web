package router

import (
	"cubeWeb/middlewares"
	"cubeWeb/model"
	"cubeWeb/utils"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

func SetPostRouter(rg *gin.RouterGroup) {
	group := rg.Group("/post")

	group.GET("/",
		getPost,
	)
	group.POST("/write/",
		middlewares.AuthRequired(),
		middlewares.CheckReqBody[writeReqBody](),
		middlewares.CheckSessionId(),
		postWrite,
	)
	group.GET("/search/", getSearch)
	group.POST("/like/",
		middlewares.AuthRequired(),
		middlewares.CheckReqBody[likeReqBody](),
		middlewares.CheckSessionId(),
		postLike,
	)
	group.POST("/unlike/",
		middlewares.AuthRequired(),
		middlewares.CheckReqBody[likeReqBody](),
		middlewares.CheckSessionId(),
		postUnlike,
	)
}

type postResData struct {
	PostId     int    `json:"postId,omitempty"`
	UserId     int    `json:"userId,omitempty"`
	Content    string `json:"content,omitempty"`
	CreateDt   int64  `json:"createDate,omitempty"`
	LikeCnt    int    `json:"likeCount"`
	UserNm     string `json:"userName,omitempty"`
	UserNickNm string `json:"userNickname,omitempty"`
}

func getPost(c *gin.Context) {
	postIdQur := c.Query("postId")
	postId, err := strconv.Atoi(postIdQur)
	if err != nil {
		utils.ResError(c, http.StatusBadRequest, 1, "Invaild postId query")
		return
	}

	var dbQuery model.Post
	if err = model.DB.Model(&model.Post{}).Select(
		"posts.*",
		"users.user_nm", "users.user_nick_nm",
	).Joins(
		"left join users on posts.user_id = users.user_id",
	).Where(
		"posts.post_id = ?", postId,
	).Find(&dbQuery).Error; err != nil {
		utils.InternalError(c, err)
		return
	}

	if dbQuery.PostId == 0 {
		utils.ResError(c, http.StatusNotFound, 2, "Post not found")
		return
	}

	utils.ResOK(c, postResData{
		dbQuery.PostId,
		dbQuery.UserId,
		dbQuery.Content,
		dbQuery.CreateDt.Unix(),
		dbQuery.LikeCnt,
		dbQuery.UserNm,
		dbQuery.UserNickNm,
	})
}

type writeReqBody struct {
	Content string `json:"content"`
}

func postWrite(c *gin.Context) {
	var reqBody writeReqBody
	utils.GetBodyJSON(c, &reqBody)

	userId, _ := utils.GetUserId(c)

	if err := model.DB.Create(
		model.Post{
			UserId:  userId,
			Content: reqBody.Content,
		},
	).Error; err != nil {
		utils.InternalError(c, err)
		return
	}

	utils.ResOK(c, nil)
}

type searchResData struct {
	PostId     int    `gorm:"not null;primaryKey;AUTO_INCREMENT" json:"postId,omitempty"`
	UserId     int    `gorm:"not null" json:"userId,omitempty"`
	Content    string `gorm:"not null" json:"content,omitempty"`
	CreateDt   int64  `gorm:"not null;datetime:timestamp;autoCreateTime" json:"createDate,omitempty"`
	LikeCnt    int    `gorm:"not null;default:0" json:"likeCount"`
	UserNm     string `json:"userName,omitempty"`
	UserNickNm string `json:"userNickname,omitempty"`
}

func getSearch(c *gin.Context) {
	userId := c.DefaultQuery("userId", "-1")
	dateFromQur := c.DefaultQuery("dateFrom", "-1")
	lengthQur := c.DefaultQuery("length", "10")
	keywordsQur := c.DefaultQuery("keywords", "")

	//NOTE - Query Keywords
	keywords := ""
	if keywordsQur != "" {
		keywords = "%(" + strings.ReplaceAll(keywordsQur, " ", "|") + ")%"
	}

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

	dbQuery := model.DB.Model(
		&model.Post{},
	).Select(
		"posts.*",
		"users.user_nm", "users.user_nick_nm",
	).Joins(
		"left join users on posts.user_id = users.user_id",
	).Order("posts.create_dt").Limit(
		length,
	).Where(
		"posts.create_dt < ?",
		dateFrom,
	)
	if userId != "-1" {
		dbQuery.Where("posts.user_id = ?", userId)
	}
	if keywords != "" {
		dbQuery.Where("posts.content SIMILAR TO ?", keywords)
	}

	if err = dbQuery.Find(&posts).Error; err != nil {
		utils.InternalError(c, err)
		return
	}

	res := []searchResData{}
	for _, value := range posts {
		res = append(res, searchResData{
			value.PostId,
			value.UserId,
			value.Content,
			value.CreateDt.Unix(),
			value.LikeCnt,
			value.UserNm,
			value.UserNickNm,
		})
	}

	utils.ResOK(c, res)
}

type likeReqBody struct {
	PostId int `json:"postId"`
}

func postLike(c *gin.Context) {
	var reqBody likeReqBody
	utils.GetBodyJSON(c, &reqBody)

	userId, _ := utils.GetUserId(c)

	var isExist bool

	err := model.DB.Model(
		&model.PostLike{},
	).Select("count(*) > 0").Where(
		"user_id = ? AND target_post_id = ?",
		userId,
		reqBody.PostId,
	).Find(&isExist).Error
	if err != nil {
		utils.InternalError(c, err)
		return
	}

	if isExist {
		utils.ResError(c, http.StatusBadRequest, 5, "You already liked")
	}

	if err := model.DB.Create(
		model.PostLike{
			UserId:       userId,
			TargetPostId: reqBody.PostId,
		},
	).Error; err != nil {
		utils.InternalError(c, err)
		return
	}

	utils.ResOK(c, nil)
}

func postUnlike(c *gin.Context) {
	var reqBody likeReqBody
	utils.GetBodyJSON(c, &reqBody)

	userId, _ := utils.GetUserId(c)

	var isExist bool

	err := model.DB.Model(
		&model.PostLike{},
	).Select("count(*) > 0").Where(
		"user_id = ? AND target_post_id = ?",
		userId,
		reqBody.PostId,
	).Find(&isExist).Error
	if err != nil {
		utils.InternalError(c, err)
		return
	}

	if !isExist {
		utils.ResError(c, http.StatusBadRequest, 5, "You did not liked")
		return
	}

	if err := model.DB.Delete(
		model.PostLike{
			UserId:       userId,
			TargetPostId: reqBody.PostId,
		},
	).Error; err != nil {
		utils.InternalError(c, err)
		return
	}

	utils.ResOK(c, nil)
}

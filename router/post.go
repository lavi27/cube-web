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
	PostId     int    `json:"postId"`
	UserId     int    `json:"userId"`
	Content    string `json:"content"`
	CreateDt   int64  `json:"createDate"`
	LikeCnt    int    `json:"likeCount"`
	UserNickNm string `json:"userNickname"`
	IsLiked    bool   `json:"isLiked"`
}

func getPost(c *gin.Context) {
	postIdQur := c.Query("postId")
	userId, userIdErr := utils.GetUserId(c)
	postId, err := strconv.Atoi(postIdQur)
	if err != nil {
		utils.ResError(c, http.StatusBadRequest, 1, "Invaild postId query")
		return
	}

	var dbRes model.Post
	selectArr := []string{
		"posts.post_id", "posts.content", "posts.create_dt",
		"users.user_id", "users.user_nm", "users.user_nick_nm",
		"COUNT(post_likes.*) as like_cnt",
	}
	groupArr := "posts.post_id, posts.content, posts.create_dt, users.user_nm, users.user_nick_nm, users.user_id"
	dbQuery := model.DB.Model(&model.Post{})

	if userIdErr != utils.ErrNotFound {
		selectArr = append(selectArr, "(case when post_likes.user_id = "+strconv.Itoa(userId)+" then true else false end) as is_liked")
		groupArr += ", post_likes.user_id"
	}

	err = dbQuery.Select(selectArr).Where("post_id = ?", postId).Group(groupArr).Joins(
		"left join users on posts.user_id = users.user_id left join post_likes on posts.post_id = post_likes.target_post_id",
	).Find(&dbRes).Error
	if err != nil {
		utils.InternalError(c, err)
		return
	}

	if dbRes.PostId == 0 {
		utils.ResError(c, http.StatusNotFound, 2, "Post not found")
		return
	}

	if userIdErr == utils.ErrNotFound {
		dbRes.IsLiked = false
	}

	utils.ResOK(c, postResData{
		dbRes.PostId,
		dbRes.UserId,
		dbRes.Content,
		dbRes.CreateDt.Unix(),
		dbRes.LikeCnt,
		dbRes.UserNickNm,
		dbRes.IsLiked,
	})
}

type writeReqBody struct {
	Content string `json:"content"`
}

func postWrite(c *gin.Context) {
	var reqBody writeReqBody
	utils.GetBodyJSON(c, &reqBody)

	userId, _ := utils.GetUserId(c)

	post := model.Post{
		UserId:  userId,
		Content: reqBody.Content,
	}

	if err := model.DB.Create(&post).Error; err != nil {
		utils.InternalError(c, err)
		return
	}

	utils.ResOK(c, nil)
}

func getSearch(c *gin.Context) {
	userIdQur := c.DefaultQuery("userId", "-1")
	dateFromQur := c.DefaultQuery("dateFrom", "-1")
	lengthQur := c.DefaultQuery("length", "10")
	keywordsQur := c.DefaultQuery("keywords", "")

	userId, userIdErr := utils.GetUserId(c)

	//NOTE - Query Keywords
	keywords := ""
	if keywordsQur != "" {
		keywords = "%(" + strings.ReplaceAll(keywordsQur, " ", "|") + ")%"
	}

	//NOTE - Query UserId
	if _, err := strconv.Atoi(userIdQur); err != nil {
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
		utils.ResError(c, http.StatusBadRequest, 4, "Too much length, it should be under 100")
		return
	}

	//NOTE - DB Create Query
	var posts []model.Post
	selectArr := []string{
		"posts.post_id", "posts.content", "posts.create_dt",
		"users.user_id", "users.user_nm", "users.user_nick_nm",
		"COUNT(post_likes.*) as like_cnt",
	}
	groupArr := "posts.post_id, posts.content, posts.create_dt, users.user_nm, users.user_nick_nm, users.user_id"
	dbQuery := model.DB.Model(&model.Post{})

	if userIdErr != utils.ErrNotFound {
		selectArr = append(selectArr, "(case when post_likes.user_id = "+strconv.Itoa(userId)+" then true else false end) as is_liked")
		groupArr += ", post_likes.user_id"
	}

	dbQuery.Select(selectArr).Group(groupArr).Joins(
		"left join users on posts.user_id = users.user_id left join post_likes on posts.post_id = post_likes.target_post_id",
	)
	dbQuery.Order("posts.create_dt DESC").Limit(length).Where(
		"posts.create_dt < ?", dateFrom,
	)

	if userIdQur != "-1" {
		dbQuery.Where("posts.user_id = ?", userIdQur)
	}
	if keywords != "" {
		dbQuery.Where("posts.content SIMILAR TO ?", keywords)
	}

	if err = dbQuery.Find(&posts).Error; err != nil {
		utils.InternalError(c, err)
		return
	}

	res := []postResData{}
	for _, value := range posts {
		if userIdErr == utils.ErrNotFound {
			value.IsLiked = false
		}

		res = append(res, postResData{
			value.PostId,
			value.UserId,
			value.Content,
			value.CreateDt.Unix(),
			value.LikeCnt,
			value.UserNickNm,
			value.IsLiked,
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

	err := model.DB.Model(&model.PostLike{}).Where(
		"user_id = ? AND target_post_id = ?", userId, reqBody.PostId,
	).Select("count(*) > 0").Find(&isExist).Error
	if err != nil {
		utils.InternalError(c, err)
		return
	}

	if isExist {
		utils.ResError(c, http.StatusBadRequest, 5, "You already liked")
		return
	}

	if err := model.DB.Create(&model.PostLike{
		UserId:       userId,
		TargetPostId: reqBody.PostId,
	}).Error; err != nil {
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

	err := model.DB.Model(&model.PostLike{}).Where(
		"user_id = ? AND target_post_id = ?", userId, reqBody.PostId,
	).Select("count(*) > 0").Find(&isExist).Error
	if err != nil {
		utils.InternalError(c, err)
		return
	}

	if !isExist {
		utils.ResError(c, http.StatusBadRequest, 5, "You did not liked")
		return
	}

	if err := model.DB.Where(
		"user_id = ? AND target_post_id = ?", userId, reqBody.PostId,
	).Delete(&model.PostLike{}).Error; err != nil {
		utils.InternalError(c, err)
		return
	}

	utils.ResOK(c, nil)
}

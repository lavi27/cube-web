package model

import (
	"cubeWeb/env"
	"cubeWeb/utils"
	"time"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

type Post struct {
	PostId     int       `gorm:"not null;primaryKey;AUTO_INCREMENT" json:"postId,omitempty"`
	UserId     int       `gorm:"not null" json:"userId,omitempty"`
	Content    string    `gorm:"not null" json:"content,omitempty"`
	CreateDt   time.Time `gorm:"not null;datetime:timestamp;autoCreateTime" json:"createDate,omitempty"`
	LikeCnt    int       `gorm:"not null;default:0" json:"likeCount,omitempty"`
	UserNm     string    `gorm:"-" json:"userName,omitempty"`
	UserNickNm string    `gorm:"-" json:"userNickname,omitempty"`
}

type User struct {
	UserId       int       `gorm:"not null;primaryKey;AUTO_INCREMENT" json:"userId,omitempty"`
	UserNm       string    `gorm:"not null" json:"userName,omitempty"`
	UserNickNm   string    `json:"userNickname,omitempty"`
	UserPw       string    `gorm:"not null;type:character(80)" json:"userPassword,omitempty"`
	CreateDt     time.Time `gorm:"not null;datetime:timestamp;autoCreateTime" json:"createDate,omitempty"`
	FollowerCnt  int       `gorm:"not null;default:0" json:"followerCount"`
	FollowingCnt int       `gorm:"not null;default:0" json:"followingCount"`
	PostCnt      int       `gorm:"-" json:"postCount"`
}

var DB *gorm.DB

func ConnectDB() {
	dsn := "host=" + env.DBIP + " user=" + env.DBUser + " password=" + env.DBPassword + " dbname=" + env.DBName + " port=" + env.DBPort + " sslmode=disable"

	config := gorm.Config{
		TranslateError: true,
		Logger:         &utils.GormLogger{},
	}

	db, err := gorm.Open(postgres.Open(dsn), &config)
	if err != nil {
		panic("DB 연결에 실패하였습니다.")
	}

	err = db.AutoMigrate(&Post{}, &User{})
	if err != nil {
		panic("DB 연결에 실패하였습니다.")
	}

	DB = db
}

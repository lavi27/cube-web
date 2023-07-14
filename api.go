package main

import (
  "net/http"
  "fmt"
  "os"
  "log"
  "database/sql" 
  // "router/home"
  
  "github.com/joho/godotenv"
  "github.com/gin-gonic/gin"
  _ "github.com/go-sql-driver/mysql"
)

type Post struct {
  content       string
  authorId      string
  timestamp     string
  likeCount     int
  commentCount  int
}

func init() {
  err := godotenv.Load(".env")
    
  if err != nil {
      log.Fatal("Error loading .env file")
  }
  
  // os.Getenv("DBNAME")
}

// const (
//   DB_USER = ""
//   DB_PASS = ""
//   DB_HOST = ""
// )

func setupRouter() *gin.Engine {
	r := gin.Default()

  home := r.Group("/")
  {
    home.GET("/", func(c *gin.Context) {
      var content string
      var authorId string
      var timestamp string
      var likeCount int
      var commentCount int
      var posts []Post
      pageId := c.Query("pageId")

      rows, err := db.Query("SELECT * FROM post WHERE id <= ?", pageId)
      if err != nil {
          log.Fatal(err)
      }
      defer rows.Close()
 
      for rows.Next() {
          err := rows.Scan(&content, &authorId, &timestamp, &likeCount, &commentCount)
          if err != nil {
              log.Fatal(err)
          }
          posts = append(posts, Post{content: content, authorId: authorId, timestamp: timestamp, likeCount: likeCount, commentCount: commentCount}) 
      }

      //posts
      c.JSON(http.StatusOK, gin.H{"posts": posts})
    })
  }

  signin := r.Group("/signin")
  {
    signin.POST("/", func(c *gin.Context) {
      c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
    })
  }

  signup := r.Group("/signup")
  {
    signup.POST("/", func(c *gin.Context) {
      c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
    })
  }

	return r
}

func main() {
  db, err := sql.Open("mysql", DB_URL)
  if err != nil {
      log.Fatal(err)
  }
  defer db.Close()

  r := setupRouter()

  r.Run() // Listen and Server in 0.0.0.0:8080
}
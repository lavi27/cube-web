package main

import (
  "net/http"
  "fmt"
  "os"
  "log"
  
  "github.com/joho/godotenv"
  "github.com/gin-gonic/gin"
)

func init() {
  err := godotenv.Load(".env")
    
  if err != nil {
      log.Fatal("Error loading .env file")
  }
  
  // os.Getenv("DBNAME")
}
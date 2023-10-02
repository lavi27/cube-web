package env

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

var (
	ClientIP   string
	DBIP       string
	DBPort     string
	DBName     string
	DBUser     string
	DBPassword string
	Enviroment string
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		log.Fatal("Error loading .env file")
	}

	ClientIP = os.Getenv("CLIENT_IP")
	DBIP = os.Getenv("DB_IP")
	DBPort = os.Getenv("DB_PORT")
	DBName = os.Getenv("DB_NAME")
	DBUser = os.Getenv("DB_USER")
	DBPassword = os.Getenv("DB_PASS")
	Enviroment = os.Getenv("ENVIROMENT")
}

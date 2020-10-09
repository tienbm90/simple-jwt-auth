package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/simple-jwt-auth/models"
	"github.com/simple-jwt-auth/servers"
	"log"
	"os"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func LoadEnv() models.Enviroment {
	var err error
	err = godotenv.Load()

	if err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	} else {
		fmt.Println("We are getting the env values")
	}
	redis_host := os.Getenv("REDIS_HOST")
	redis_port := os.Getenv("REDIS_PORT")
	redis_password := os.Getenv("REDIS_PASSWORD")
	redis := models.RedisConf{redis_host, redis_port, redis_password}
	sql_username := os.Getenv("SQL_USERNAME")
	sql_password := os.Getenv("SQL_PASSWORD")
	sql_url := os.Getenv("SQL_URL")
	sql := models.SqlConf{sql_username, sql_password, sql_url}

	port := os.Getenv("APP_PORT")
	env := models.Enviroment{redis, sql, port}
	return env
}

var AuthServer servers.Server

func main() {
	//servers.Run()
	Run()
	log.Println("Server exiting")
}

func Run() {
	AuthServer = servers.Server{}

	env := LoadEnv()
	AuthServer.Initialize(env)
	AuthServer.Run()
	//AuthServer.CheckEnforcer(env)
}

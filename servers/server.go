package core

import (
	"fmt"
	"github.com/casbin/casbin/persist/file-adapter"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/joho/godotenv"
	"github.com/simple-jwt-auth/auth"
	"log"
	"net/http"
	"os"
)

type Server struct {
	Router      *gin.Engine
	FileAdapter *fileadapter.Adapter
	RedisCli    *redis.Client
	rd          auth.AuthInterface
	tk          auth.TokenInterface
}

var server Server

func (server *Server) Initialize(redis_host, redis_port, redis_password string) {
	server.Router = gin.Default()
	server.RedisCli = NewRedisDB(redis_host, redis_port, redis_password)
	server.FileAdapter = fileadapter.NewAdapter("config/basic_policy.csv")
	//init route
	server.rd = auth.NewAuthService(server.RedisCli)
	server.tk = auth.NewTokenService()
	server.InitializeRoutes()
}

func NewRedisDB(host, port, password string) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       0,
	})
	return redisClient
}

func (server *Server) Run(addr string) {
	fmt.Printf("Listen on port %s \n", addr)
	log.Fatal(http.ListenAndServe(addr, server.Router))
}

func Run() {
	var server = Server{}
	var err error
	err = godotenv.Load()
	if err != nil {
		log.Fatalf("Error getting env, not comming through %v", err)
	} else {
		fmt.Println("We are getting the env values")
	}

	appAddr := ":" + os.Getenv("PORT")

	//redis details
	redis_host := os.Getenv("REDIS_HOST")
	redis_port := os.Getenv("REDIS_PORT")
	redis_password := os.Getenv("REDIS_PASSWORD")

	server.Initialize(redis_host, redis_port, redis_password)

	server.Run(":" + appAddr)

}

func (server *Server) Close() {
	//close DB
	if server.RedisCli != nil {
		if err := server.RedisCli.Close(); err != nil {
			log.Fatal(err)
		}

		server.RedisCli = nil
	}
}

package main

import (
	"context"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/gophercon-jwt-repo/auth"
	"github.com/gophercon-jwt-repo/handlers"
	"github.com/gophercon-jwt-repo/middleware"
	"github.com/joho/godotenv"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"
	//"github.com/casbin/casbin"
	fileadapter "github.com/casbin/casbin/persist/file-adapter"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}

func NewRedisDB(host, port, password string) *redis.Client {
	redisClient := redis.NewClient(&redis.Options{
		Addr:     host + ":" + port,
		Password: password,
		DB:       0,
	})
	return redisClient
}

func main() {

	appAddr := ":" + os.Getenv("PORT")

	//redis details
	redis_host := os.Getenv("REDIS_HOST")
	redis_port := os.Getenv("REDIS_PORT")
	redis_password := os.Getenv("REDIS_PASSWORD")

	redisClient := NewRedisDB(redis_host, redis_port, redis_password)
	adapter := fileadapter.NewAdapter("config/basic_policy.csv")

	var rd = auth.NewAuth(redisClient)
	var tk = auth.NewToken()
	var service = handlers.NewProfile(rd, tk)

	var router = gin.Default()

	router.POST("/login", service.Login)

	resource := router.Group("/api")

	resource.Use(middleware.TokenAuthMiddleware())
	{
		resource.POST("/todo",  middleware.Authorize("resource", "write", adapter), service.CreateTodo)
		resource.GET("/todo", middleware.Authorize("resource", "read", adapter), service.GetTodo)
		resource.POST("/logout", service.Logout)
		resource.POST("/refresh", service.Refresh)
	}


	//router.POST("/todo", middleware.TokenAuthMiddleware(), service.CreateTodo)
	//router.POST("/logout", middleware.TokenAuthMiddleware(), service.Logout)
	//router.POST("/refresh", service.Refresh)

	srv := &http.Server{
		Addr:              appAddr,
		Handler:           router,
	}
	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()
	//Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	log.Println("Shutdown Server ...")

	ctx, cancel := context.WithTimeout(context.Background(), 10 * time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server Shutdown:", err)
	}
	log.Println("Server exiting")
}

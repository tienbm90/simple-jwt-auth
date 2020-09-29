package main

import (
	"github.com/joho/godotenv"
	"github.com/simple-jwt-auth/servers"
	"log"
)

func init() {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}
}


func main() {
	servers.Run()
	log.Println("Server exiting")
}

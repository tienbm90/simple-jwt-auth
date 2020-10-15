package main

import (
	"fmt"
	"github.com/joho/godotenv"
	"github.com/simple-jwt-auth/models"
	"github.com/simple-jwt-auth/servers"
	"log"
	"os"
	"strconv"
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
	redisHost := os.Getenv("REDIS_HOST")
	redisPort := os.Getenv("REDIS_PORT")
	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisUsername := os.Getenv("REDIS_USERNAME")
	var redis = models.RedisConf{Host: redisHost, Port: redisPort, Username: redisUsername, Password: redisPassword}
	dbDriver := os.Getenv("DB_DRIVER")
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbDatabase := os.Getenv("DB_DATABASE")
	sqlConf := models.SqlConf{Driver: dbDriver, Username: dbUsername, Passord: dbPassword, Host: dbHost, Port: dbPort, Database: dbDatabase}

	googleClientId := os.Getenv("GOOGLE_CLIENT_ID")
	googleProjectId := os.Getenv("GOOGLE_PROJECT_ID")
	googleAuthUri := os.Getenv("GOOGLE_AUTH_URI")
	googleTokenUri := os.Getenv("GOOGLE_TOKEN_URI")
	googleCertUri := os.Getenv("GOOGLE_AUTH_PROVIDER_X509_CERT_URL")
	googleClientSecret := os.Getenv("GOOGLE_CLIENT_SECRET")
	googleRedirectUrl := os.Getenv("GOOGLE_REDIRECT_URL")
	gConf := models.Google{
		ClientID:                googleClientId,
		ProjectID:               googleProjectId,
		AuthUri:                 googleAuthUri,
		TokenUri:                googleTokenUri,
		AuthProviderX509CertUri: googleCertUri,
		ClientSecret:            googleClientSecret,
		RedirectUrl:             googleRedirectUrl,
	}

	githubClientId := os.Getenv("GITHUB_CLIENT_ID")
	githubClientSecret := os.Getenv("GITHUB_CLIENT_SECRET")
	githubRedirectUrl := os.Getenv("GITHUB_REDIRECT_URL")

	githubConf := models.Github{
		ClientID:     githubClientId,
		ClientSecret: githubClientSecret,
		RedirectUrl:  githubRedirectUrl,
	}

	facebookClientId := os.Getenv("FACEBOOK_CLIENT_ID")
	facebookClientSecret := os.Getenv("FACEBOOK_CLIENT_SECRET")
	facebookRedirecUrl := os.Getenv("FACEBOOK_REDIRECT_URL")

	facebookConf := models.Facebook{
		ClientID:     facebookClientId,
		ClientSecret: facebookClientSecret,
		RedirectUrl:  facebookRedirecUrl,
	}

	port := os.Getenv("APP_PORT")

	watcherEnable, err := strconv.ParseBool(os.Getenv("CASBIN_WATCHER_ENABLE"))

	if err != nil {
		log.Fatal(fmt.Sprintf("Can't parse enviroment: %s", err.Error()))
	}
	env := models.Enviroment{
		RedisConfig:         redis,
		SqlConfig:           sqlConf,
		Port:                port,
		CasbinWatcherEnable: watcherEnable,
		FacebookConf:        facebookConf,
		GithubConf:          githubConf,
		GoogleConf:          gConf}
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

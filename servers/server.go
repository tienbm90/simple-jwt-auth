package servers

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/simple-jwt-auth/auth"
	"github.com/simple-jwt-auth/models"
	"github.com/simple-jwt-auth/utils"
	"log"
	"net/http"
)

type Server struct {
	Router     *gin.Engine
	Enforcer   *casbin.Enforcer
	RedisCli   *redis.Client
	RD         auth.AuthInterface
	TK         auth.TokenInterface
	enviroment models.Enviroment
}

func (server *Server) Initialize(env models.Enviroment) {
	server.enviroment = env
	server.Router = gin.Default()
	server.RedisCli = utils.NewRedisDB(server.enviroment.RedisConfig.Host, server.enviroment.RedisConfig.Port, server.enviroment.RedisConfig.Password)
	dataSource := fmt.Sprintf("%s:%s@tcp(%s)/", server.enviroment.SqlConfig.Username, server.enviroment.SqlConfig.Passord, server.enviroment.SqlConfig.Url)
	server.Enforcer = auth.NewCasbinEnforcer(dataSource)

	//init route
	server.RD = auth.NewAuthService(server.RedisCli)
	server.InitializeRoutes()
}

func (server *Server) Run() {
	fmt.Printf("Listen on port %s \n", server.enviroment.Port)
	//log.Fatal(http.ListenAndServe(server.enviroment.Port, server.Router))
	log.Fatal(http.ListenAndServe(":8081", server.Router))
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

func (server *Server) CheckEnforcer(env models.Enviroment) {
	server.enviroment = env
	server.Router = gin.Default()
	server.RedisCli = utils.NewRedisDB(server.enviroment.RedisConfig.Host, server.enviroment.RedisConfig.Port, server.enviroment.RedisConfig.Password)
	dataSource := fmt.Sprintf("%s:%s@tcp(%s)/", server.enviroment.SqlConfig.Username, server.enviroment.SqlConfig.Passord, server.enviroment.SqlConfig.Url)
	server.Enforcer = auth.NewCasbinEnforcer(dataSource)

	ok, err := server.Enforcer.Enforce("admin", "/auth/policy", "GET")
	if err != nil {
		log.Fatal(fmt.Sprintf("Error: %s", err.Error()))
	}

	ok, err = server.Enforcer.Enforce("admin", "/auth/policy/1", "GET")
	if err != nil {
		log.Fatal(fmt.Sprintf("Error: %s", err.Error()))
	}

	fmt.Println(ok)
}

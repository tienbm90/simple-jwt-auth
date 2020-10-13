package servers

import (
	"fmt"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v7"
	"github.com/simple-jwt-auth/auth"
	"github.com/simple-jwt-auth/models"
	"github.com/simple-jwt-auth/seed"
	"github.com/simple-jwt-auth/utils"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

type Server struct {
	Router     *gin.Engine
	Enforcer   *casbin.Enforcer
	RedisCli   *redis.Client
	RD         auth.AuthInterface
	TK         auth.JwtTokenInterface
	enviroment models.Enviroment
	DB         *gorm.DB
}

func (server *Server) Initialize(env models.Enviroment) {
	server.enviroment = env
	server.Router = gin.Default()
	server.RedisCli = utils.NewRedisDB(
		server.enviroment.RedisConfig.Host,
		server.enviroment.RedisConfig.Port,
		server.enviroment.RedisConfig.Password)

	server.InitDB("mysql",
		env.SqlConfig.Username,
		env.SqlConfig.Passord,
		env.SqlConfig.Host,
		env.SqlConfig.Port,
		env.SqlConfig.Database)

	server.Enforcer = auth.NewCasbinEnforcerFromDB(server.DB)
	//init route
	server.RD = auth.NewAuthService(server.RedisCli)
	seed.Load(server.DB)
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

//func (server *Server) CheckEnforcer(env models.Enviroment) {
//	server.enviroment = env
//	server.Router = gin.Default()
//	server.RedisCli = utils.NewRedisDB(server.enviroment.RedisConfig.Host, server.enviroment.RedisConfig.Port, server.enviroment.RedisConfig.Password)
//	dataSource := fmt.Sprintf("%s:%s@tcp(%s)/", server.enviroment.SqlConfig.Username, server.enviroment.SqlConfig.Passord, server.enviroment.SqlConfig.Address)
//	server.Enforcer = auth.NewCasbinEnforcer(dataSource)
//
//	ok, err := server.Enforcer.Enforce("admin", "/auth/policy", "GET")
//	if err != nil {
//		log.Fatal(fmt.Sprintf("Error: %s", err.Error()))
//	}
//
//	ok, err = server.Enforcer.Enforce("admin", "/auth/policy/1", "GET")
//	if err != nil {
//		log.Fatal(fmt.Sprintf("Error: %s", err.Error()))
//	}
//
//	fmt.Println(ok)
//}

func (server *Server) InitDB(Dbdriver, DbUser, DbPassword, DbHost, DbPort, DbName string) {
	var err error
	if Dbdriver == "mysql" {
		DBURL := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=True&loc=Local", DbUser, DbPassword, DbHost, DbPort, DbName)

		server.DB, err = gorm.Open(mysql.Open(DBURL), &gorm.Config{})
		if err != nil || server.DB == nil {
			for i := 1; i <= 12; i++ {
				fmt.Printf("gorm.Open(%s, %s) %d\n", Dbdriver, DBURL, i)
				server.DB, err = gorm.Open(mysql.Open(DBURL), &gorm.Config{})

				if server.DB != nil && err == nil {
					break
				} else {
					time.Sleep(5 * time.Second)
				}
			}

			if err != nil || server.DB == nil {
				fmt.Println(err)
				log.Fatal(err)
			}
		} else {
			fmt.Printf("We are connected to the %s database", Dbdriver)
		}
	}
	if Dbdriver == "postgres" {
		DBURL := fmt.Sprintf("host=%s port=%s user=%s dbname=%s sslmode=disable password=%s", DbHost, DbPort, DbUser, DbName, DbPassword)
		server.DB, err = gorm.Open(postgres.Open(DBURL), &gorm.Config{})
		if err != nil {

			for i := 1; i <= 12; i++ {
				fmt.Printf("Retry connect to gorm.Open(%s, %s) %d\n", Dbdriver, DBURL, i)
				server.DB, err = gorm.Open(postgres.Open(DBURL), &gorm.Config{})

				if server.DB != nil && err == nil {
					break
				} else {
					time.Sleep(5 * time.Second)
				}
			}

			if err != nil {
				fmt.Printf("Cannot connect to %s database", Dbdriver)
				log.Fatal("This is the error:", err)
				fmt.Println(err)
			}
		} else {
			fmt.Printf("We are connected to the %s database", Dbdriver)
		}

	}
	server.DB.Debug().AutoMigrate(&models.User{}, &models.Policy{}, &models.UserRole{}, &models.Todo{}) //database migration
}

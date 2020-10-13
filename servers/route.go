package servers

import (
	"github.com/gin-gonic/gin"
	"github.com/simple-jwt-auth/api"
	"github.com/simple-jwt-auth/middleware"
	"github.com/simple-jwt-auth/models"
	"golang.org/x/oauth2"
)

func (server *Server) InitializeRoutes() {
	casbinService := api.NewCasbinService(server.Enforcer)
	userRepos := models.ProvideUserRepository(server.DB)
	jwtApi := api.CreateJwtApi(&userRepos)

	googleConf := oauth2.Config{
		ClientID:     server.enviroment.GoogleConf.ClientID,
		ClientSecret: server.enviroment.GoogleConf.ClientSecret,
		Endpoint:     oauth2.Endpoint{},
		RedirectURL:  server.enviroment.GoogleConf.RedirectUrl,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email", // You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
		},
	}

	googleApi := api.GoogleAPI{
		Config:   &googleConf,
		UserRepo: &userRepos,
	}
	// jwt api
	server.Router.POST("/login/token", jwtApi.JwtLogin)
	jwt := server.Router.Group("/api")
	jwt.Use(gin.Logger())
	jwt.Use(gin.Recovery())
	jwt.Use(middleware.TokenAuthMiddleware())
	{
		jwt.POST("/auth/policy", middleware.AuthorizeJwtToken("/auth/policy", "POST", server.Enforcer), casbinService.CreatePolicy)
		jwt.GET("/auth/policy", middleware.AuthorizeJwtToken("/auth/policy", "GET", server.Enforcer), casbinService.ListPolicy)
		jwt.POST("/auth/grouppolicy", middleware.AuthorizeJwtToken("/auth/grouppolicy", "POST", server.Enforcer), casbinService.CreateGroupPolicy)
		jwt.GET("/auth/grouppolicy", middleware.AuthorizeJwtToken("/auth/grouppolicy", "GET", server.Enforcer), casbinService.ListGroupPolicies)
		jwt.POST("/api/todo", middleware.AuthorizeJwtToken("resource", "write", server.Enforcer), api.CreateTodo)
		jwt.GET("/api/todo", middleware.AuthorizeJwtToken("resource", "read", server.Enforcer), api.GetTodo)
		jwt.POST("/logout", jwtApi.JwtLogout)
		jwt.POST("/refresh", jwtApi.JwtRefresh)
	}

	// openid api
	server.Router.POST("/login/oauth", googleApi.LoginHandler)
	googleOauth := server.Router.Group("/oauth")
	googleOauth.Use(gin.Logger())
	googleOauth.Use(gin.Recovery())
	googleOauth.Use(middleware.AuthorizeOpenIdRequest())
	{
		jwt.POST("/todo", middleware.AuthorizeJwtToken("/auth/policy", "POST", server.Enforcer), casbinService.CreatePolicy)
		jwt.GET("/todo", middleware.AuthorizeJwtToken("/auth/policy", "GET", server.Enforcer), casbinService.ListPolicy)

	}
}

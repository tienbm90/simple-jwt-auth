package servers

import (
	"fmt"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/simple-jwt-auth/api"
	"github.com/simple-jwt-auth/ginserver"
	"github.com/simple-jwt-auth/middleware"
	"github.com/simple-jwt-auth/models"
	"github.com/simple-jwt-auth/utils"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	"golang.org/x/oauth2/google"

	//simple_models "github.com/simple-jwt-auth/models"
	//"golang.org/x/oauth2/facebook"
	githuboauth "golang.org/x/oauth2/github"
	"log"
)

func (server *Server) InitializeRoutes() {

	////init casbinservice
	casbinService := api.NewCasbinService(server.Enforcer)
	userRepos := models.ProvideUserRepository(server.DB)
	oauthClientRepo := models.OauthClientRepository{DB: server.DB}

	token, err := middleware.RandToken(64)
	if err != nil {
		log.Fatal("unable to generate random token: ", err)
	}
	store := sessions.NewCookieStore([]byte(token))

	store.Options(sessions.Options{
		Path:   "/",
		MaxAge: 86400 * 7,
	})

	userApi := api.ProvideUserAPI(api.UserService{UserRepository: userRepos})

	//create jwt api
	jwtApi := api.CreateJwtApi(&userRepos)

	//create google api
	googleConf := oauth2.Config{
		ClientID:     server.enviroment.GoogleConf.ClientID,
		ClientSecret: server.enviroment.GoogleConf.ClientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  server.enviroment.GoogleConf.RedirectUrl,
		Scopes: []string{
			//"https://www.googleapis.com/oauth2/userinfo.email", // You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
			"https://www.googleapis.com/auth/userinfo.email", // You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
		},
	}

	googleApi := api.GoogleAPI{
		Config:   &googleConf,
		UserRepo: &userRepos,
	}

	// create facebook api
	facebookConf := oauth2.Config{
		ClientID:     server.enviroment.FacebookConf.ClientID,
		ClientSecret: server.enviroment.FacebookConf.ClientSecret,
		Endpoint:     facebook.Endpoint,
		RedirectURL:  server.enviroment.FacebookConf.RedirectUrl,
		Scopes: []string{
			"email",
			"public_profile",
			//"user_link",
			//"user_localtion",
		},
	}
	facebookApi := api.FacebookAPI{
		Config:   &facebookConf,
		UserRepo: &userRepos,
	}

	// create github api
	githubConf := &oauth2.Config{
		ClientID:     server.enviroment.GithubConf.ClientID,
		ClientSecret: server.enviroment.GithubConf.ClientSecret,
		Scopes:       []string{"user"},
		Endpoint:     githuboauth.Endpoint,
	}

	githubApi := api.GithubAPI{
		Config:   githubConf,
		UserRepo: &userRepos,
	}

	// init api

	server.Router.Use(gin.Logger())
	server.Router.Use(gin.Recovery())
	server.Router.Static("/css", "./static/css")
	server.Router.Static("/img", "./static/img")
	server.Router.LoadHTMLGlob("./static/templates/**/*")

	//// jwt
	jwt := server.Router.Group("/jwt")
	jwt.POST("/login", jwtApi.JwtLogin)
	jwt.Use(middleware.TokenAuthMiddleware())
	{
		jwt.POST("/oauth2/policy", middleware.AuthorizeJwtToken("/jwt/oauth2/policy", "POST", server.Enforcer), casbinService.CreatePolicy)
		jwt.GET("/oauth2/policy", middleware.AuthorizeJwtToken("/jwt/oauth2/policy", "GET", server.Enforcer), casbinService.ListPolicy)
		jwt.DELETE("/oauth2/policy", middleware.AuthorizeJwtToken("/jwt/oauth2/policy", "DELETE", server.Enforcer), casbinService.DeletePolicy)
		jwt.POST("/oauth2/grouppolicy", middleware.AuthorizeJwtToken("/jwt/oauth2/grouppolicy", "POST", server.Enforcer), casbinService.CreateGroupPolicy)
		jwt.GET("/oauth2/grouppolicy", middleware.AuthorizeJwtToken("/jwt/oauth2/grouppolicy", "GET", server.Enforcer), casbinService.ListGroupPolicies)
		//jwt.POST("/todo", middleware.AuthorizeJwtToken("resource", "write", server.Enforcer), api.CreateTodo)
		//jwt.GET("/todo", middleware.AuthorizeJwtToken("resource", "read", server.Enforcer), api.GetTodo)
		jwt.POST("/logout", jwtApi.JwtLogout)
		jwt.POST("/refresh", jwtApi.JwtRefresh)
	}
	//
	//// init route for google api
	googleOauth := server.Router.Group("/oauth/google")
	googleOauth.Use(sessions.Sessions("goquestsession", store))
	googleOauth.GET("/", googleApi.IndexHandler)

	googleOauth.GET("/login", googleApi.LoginHandler)
	googleOauth.GET("/auth", googleApi.AuthHandler)
	googleOauth.Use(middleware.AuthorizeOpenIdRequest())
	{
		googleOauth.GET("/field", googleApi.FieldHandler)
		googleOauth.GET("/test", googleApi.TestHandler)
	}
	//
	////init route for github api
	githubOauth := server.Router.Group("/oauth/github")
	githubOauth.Use(sessions.Sessions("goquestsession", store))
	githubOauth.GET("/", githubApi.IndexHandler)

	githubOauth.GET("/login", githubApi.LoginHandler)
	githubOauth.GET("/auth", githubApi.AuthHandler)
	githubOauth.Use(middleware.AuthorizeOpenIdRequest())
	{
		githubOauth.GET("/field", githubApi.FieldHandler)
		githubOauth.GET("/test", githubApi.TestHandler)
	}
	//
	////init route for facebook api
	facebookOauth := server.Router.Group("/oauth/facebook")
	facebookOauth.Use(sessions.Sessions("goquestsession", store))
	facebookOauth.GET("/", facebookApi.IndexHandler)

	facebookOauth.GET("/login", facebookApi.LoginHandler)
	facebookOauth.GET("/auth", facebookApi.AuthHandler)
	facebookOauth.Use(middleware.AuthorizeOpenIdRequest())
	{
		facebookOauth.GET("/field", facebookApi.FieldHandler)
		facebookOauth.GET("/test", facebookApi.TestHandler)
	}

	//init route for oauth2 api

	server.Router.Use(sessions.Sessions("goquestsession", store))
	//oauth2_api := api.ProviderOauth2API()
	oauth2_api := api.InitOauth2API(&oauthClientRepo, &userRepos)
	oauth2 := server.Router.Group(fmt.Sprintf("/%s", utils.OAUTH2_PREFIX))
	{
		oauth2.GET("/login", oauth2_api.ShowLoginPage)
		oauth2.POST("/login", oauth2_api.HandleLogin)
		oauth2.GET("/authorize", oauth2_api.Authorize)
		oauth2.POST("/authorize", oauth2_api.Authorize)
		oauth2.GET("/auth", oauth2_api.Authenticate)
		oauth2.POST("/auth", oauth2_api.Authenticate)
		//oauth2.GET("/token", oauth2_api.HandleTokenRequest)
		oauth2.GET("/token", ginserver.HandleAuthorizeRequest)
		//oauth2.POST("/token", oauth2_api.HandleTokenRequest)
		oauth2.POST("/token", ginserver.HandleTokenRequest)
		oauth2.GET("/test", oauth2_api.Test)
	}

	oauth2.Use(middleware.AuthorizeOpenIdRequest())
	{
		oauth2.POST("/userinfo", userApi.UserInfo)
		oauth2.GET("/userinfo", userApi.UserInfo)
	}
	server.Router.GET("/userinfo", userApi.UserInfo)
}

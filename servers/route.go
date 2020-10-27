package servers

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	ginserver "github.com/go-oauth2/gin-server"
	"github.com/simple-jwt-auth/api"
	"github.com/simple-jwt-auth/middleware"
	simple_models "github.com/simple-jwt-auth/models"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/facebook"
	githuboauth "golang.org/x/oauth2/github"
	"golang.org/x/oauth2/google"
	"log"
	"net/http"
)

func (server *Server) InitializeRoutes() {
	casbinService := api.NewCasbinService(server.Enforcer)
	userRepos := simple_models.ProvideUserRepository(server.DB)
	jwtApi := api.CreateJwtApi(&userRepos)
	token, err := middleware.RandToken(64)
	if err != nil {
		log.Fatal("unable to generate random token: ", err)
	}
	store := sessions.NewCookieStore([]byte(token))
	store.Options(sessions.Options{
		Path:   "/",
		MaxAge: 86400 * 7,
	})
	googleConf := oauth2.Config{
		ClientID:     server.enviroment.GoogleConf.ClientID,
		ClientSecret: server.enviroment.GoogleConf.ClientSecret,
		Endpoint:     google.Endpoint,
		RedirectURL:  server.enviroment.GoogleConf.RedirectUrl,
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email", // You have to select your own scope from here -> https://developers.google.com/identity/protocols/googlescopes#google_sign-in
		},
	}

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
	googleApi := api.GoogleAPI{
		Config:   &googleConf,
		UserRepo: &userRepos,
	}

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

	facebookApi := api.FacebookAPI{
		Config:   &facebookConf,
		UserRepo: &userRepos,
	}
	// jwt api
	server.Router.POST("/jwt/login", jwtApi.JwtLogin)
	server.Router.Use(gin.Logger())
	server.Router.Use(gin.Recovery())
	server.Router.Static("/css", "./static/css")
	server.Router.Static("/img", "./static/img")
	server.Router.LoadHTMLGlob("./static/templates/*")

	// jwt
	jwt := server.Router.Group("/jwt")
	jwt.Use(middleware.TokenAuthMiddleware())
	{
		jwt.POST("/auth/policy", middleware.AuthorizeJwtToken("/jwt/auth/policy", "POST", server.Enforcer), casbinService.CreatePolicy)
		jwt.GET("/auth/policy", middleware.AuthorizeJwtToken("/jwt/auth/policy", "GET", server.Enforcer), casbinService.ListPolicy)
		jwt.DELETE("/auth/policy", middleware.AuthorizeJwtToken("/jwt/auth/policy", "DELETE", server.Enforcer), casbinService.DeletePolicy)
		jwt.POST("/auth/grouppolicy", middleware.AuthorizeJwtToken("/jwt/auth/grouppolicy", "POST", server.Enforcer), casbinService.CreateGroupPolicy)
		jwt.GET("/auth/grouppolicy", middleware.AuthorizeJwtToken("/jwt/auth/grouppolicy", "GET", server.Enforcer), casbinService.ListGroupPolicies)
		//jwt.POST("/todo", middleware.AuthorizeJwtToken("resource", "write", server.Enforcer), api.CreateTodo)
		//jwt.GET("/todo", middleware.AuthorizeJwtToken("resource", "read", server.Enforcer), api.GetTodo)
		jwt.POST("/logout", jwtApi.JwtLogout)
		jwt.POST("/refresh", jwtApi.JwtRefresh)
	}

	// googleid api
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

	//github api
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

	//facebook api
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

	//oauth2 api
	auth := server.Router.Group("/oauth2")
	{
		auth.GET("/token", ginserver.HandleTokenRequest)
	}

	api := server.Router.Group("/api")
	{
		api.Use(ginserver.HandleTokenVerify())
		api.GET("/test", func(c *gin.Context) {
			ti, exists := c.Get(ginserver.DefaultConfig.TokenKey)
			if exists {
				c.JSON(http.StatusOK, ti)
				return
			}
			c.String(http.StatusOK, "not found")
		})
	}
}

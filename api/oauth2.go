package api

import (
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	ginsvr "github.com/go-oauth2/gin-server"
	oauth2_modesl "github.com/go-oauth2/oauth2/v4/models"
	"github.com/simple-jwt-auth/models"
	"gopkg.in/oauth2.v3"
	"gopkg.in/oauth2.v3/manage"
	"gopkg.in/oauth2.v3/server"
	"gopkg.in/oauth2.v3/store"
	"net/http"
	"sync"
)

type Oauth2API struct {
	gServer *server.Server
	once    sync.Once
}

//var (
//	gServer *server.Server
//	once    sync.Once
//)

// InitServer Initialize the service
func (a *Oauth2API) InitServer(manager oauth2.Manager) *server.Server {
	a.once.Do(func() {
		a.gServer = server.NewDefaultServer(manager)
	})
	return a.gServer
}

// HandleAuthorizeRequest the authorization request handling
func (a *Oauth2API) HandleAuthorizeRequest(c *gin.Context) {
	err := a.gServer.HandleAuthorizeRequest(c.Writer, c.Request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.Abort()
}

// HandleTokenRequest token request handling
func (a *Oauth2API) HandleTokenRequest(c *gin.Context) {
	err := a.gServer.HandleTokenRequest(c.Writer, c.Request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
	c.Abort()
}
func ProviderOauth2API() Oauth2API {
	// init oauth2 server
	manager := manage.NewDefaultManager()

	// token store
	manager.MustTokenStorage(store.NewFileTokenStore("data.db"))

	// client store
	clientStore := store.NewClientStore()
	clientStore.Set("000000", &oauth2_modesl.Client{
		ID:     "000000",
		Secret: "999999",
		Domain: "http://localhost",
	})

	manager.MapClientStorage(clientStore)
	// Initialize the oauth2 service

	svr := ginsvr.InitServer(manager)
	ginsvr.SetAllowGetAccessRequest(true)
	ginsvr.SetClientInfoHandler(server.ClientFormHandler)
	api := Oauth2API{
		gServer: svr,
		once:    sync.Once{},
	}

	return api
}

func (a *Oauth2API) Login(c *gin.Context) {
	session := sessions.Default(c)
	var u models.User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	session.Set("LoggedInUserID", u.UserName)
	session.Save()

	c.Header("Location", "/auth")
	c.HTML(http.StatusOK, "oauth2/login.tmpl", gin.H{"username": u.UserName, "seen": "seen"})
}

func (a *Oauth2API) Auth(c *gin.Context) {
	session := sessions.Default(c)
	var u models.User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	logedInUserID := session.Get("LoggedInUserID")

	if logedInUserID != "" && logedInUserID != nil {
		c.Header("Location", "/login")
		c.HTML(http.StatusOK, "oauth2/login.tmpl", gin.H{"username": u.UserName, "seen": "seen"})
		return
	}

	c.Header("Location", "/auth")
	c.HTML(http.StatusOK, "oauth2/auth.tmpl", "")
}

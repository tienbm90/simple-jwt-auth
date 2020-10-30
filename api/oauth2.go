package api

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	oauth2_modesl "github.com/go-oauth2/oauth2/v4/models"
	ginsvr "github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	"github.com/go-session/session"
	"github.com/simple-jwt-auth/utils"
	"log"
	"time"

	//ginsvr "github.com/go-oauth2/gin-server"
	//"github.com/go-session/session"
	//"gopkg.in/oauth2.v3"
	//"gopkg.in/oauth2.v3/generates"
	//"gopkg.in/oauth2.v3/manage"
	//oauth2_modesl "gopkg.in/oauth2.v3/models"
	//"gopkg.in/oauth2.v3/server"
	//"gopkg.in/oauth2.v3/store"
	"net/http"
	"net/url"
	"sync"
)

type Oauth2API struct {
	gServer *ginsvr.Server
	once    sync.Once
}

func ProviderOauth2API() Oauth2API {
	// init oauth2 server

	manager := manage.NewDefaultManager()
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)

	// token store
	//manager.MustTokenStorage(store.NewFileTokenStore("data.db"))
	manager.MustTokenStorage(store.NewMemoryTokenStore())
	manager.MapAccessGenerate(generates.NewJWTAccessGenerate("", []byte("00000000"), jwt.SigningMethodHS512))

	// client store
	clientStore := store.NewClientStore()
	clientStore.Set("222222", &oauth2_modesl.Client{
		ID:     "222222",
		Secret: "22222222",
		Domain: "http://localhost:8085",
	})

	manager.MapClientStorage(clientStore)
	// Initialize the oauth2 service

	svr := ginsvr.NewServer(ginsvr.NewConfig(), manager)

	svr.SetPasswordAuthorizationHandler(func(username, password string) (userID string, err error) {
		if username == "test" && password == "test" {
			userID = "test"
		}
		return
	})

	svr.SetUserAuthorizationHandler(userAuthorizeHandler)

	svr.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})

	svr.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	api := Oauth2API{
		gServer: svr,
		once:    sync.Once{},
	}

	return api
}

// InitServer Initialize the service
//func (a *Oauth2API) InitServer(manager oauth2.Manager) *ginsvr.Server {
//	a.once.Do(func() {
//		a.gServer = ginsvr.NewDefaultServer(manager)
//	})
//	return a.gServer
//}

// HandleAuthorizeRequest the authorization request handling

func (a *Oauth2API) HandleAuthorizeRequest(c *gin.Context) {
	err := a.gServer.HandleAuthorizeRequest(c.Writer, c.Request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}

}

// HandleTokenRequest token request handling
func (a *Oauth2API) HandleTokenRequest(c *gin.Context) {
	err := a.gServer.HandleTokenRequest(c.Writer, c.Request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)
		return
	}
}

// HandleTokenRequest token request handling
func (a *Oauth2API) Test(c *gin.Context) {
	r := c.Request
	w := c.Writer
	token, err := a.gServer.ValidationBearerToken(r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	data := map[string]interface{}{
		"expires_in": int64(token.GetAccessCreateAt().Add(token.GetAccessExpiresIn()).Sub(time.Now()).Seconds()),
		"client_id":  token.GetClientID(),
		"user_id":    token.GetUserID(),
	}
	e := json.NewEncoder(w)
	e.SetIndent("", "  ")
	e.Encode(data)
}

func (a *Oauth2API) Login(c *gin.Context) {
	w := c.Writer
	store, err := session.Start(c.Request.Context(), w, c.Request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if c.Request.Method == "POST" {
		if c.Request.Form == nil {
			if err := c.Request.ParseForm(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		store.Set("LoggedInUserID", c.Request.Form.Get("username"))
		store.Save()

		w.Header().Set("Location", fmt.Sprintf("/%s/%s", utils.OAUTH2_PREFIX, "auth"))
		w.WriteHeader(http.StatusFound)
		return
	}
	//outputHTML(w, r, "static/templates/oauth2/login.html")
	c.HTML(http.StatusFound, "login.tmpl", gin.H{"link": "/oauth2/login"})
	//c.Abort()
}

func (a *Oauth2API) Authenicate(c *gin.Context) {
	//r := c.Request
	//w := c.Writer
	//store, err := session.Start(nil, w, r)
	store := sessions.Default(c)

	store.Set("state", c.Request.Form.Get("state"))
	log.Printf("Authenticate state: %s \n", store.Get("state") )


	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}
	//
	////check if user is logged in
	//if _, ok := store.Get("LoggedInUserID"); !ok {
	//	w.Header().Set("Location", fmt.Sprintf("/%s/%s", utils.OAUTH2_PREFIX, "login"))
	//	return
	//}
	//
	//v, ok := store.Get("ReturnUri")
	//if ok {
	//	log.Printf("Stored ReturnURI: %s", v)
	//}
	c.HTML(http.StatusFound, "auth.tmpl", gin.H{"link": "link"})
}

func (a *Oauth2API) Authorize(c *gin.Context) {
	r := c.Request
	w := c.Writer
	store := sessions.Default(c)

	log.Println(store.Get("state"))
	store.Set("state", "xyz")
	//sessionStore, err := session.Start(r.Context(), w, r)
	//if err != nil {
	//	http.Error(w, err.Error(), http.StatusInternalServerError)
	//	return
	//}

	var form url.Values
	form = c.Request.Form
	log.Printf("Form: %s", form)
	//sessionStore.Set("state", form.Get(""))
	//sessionStore.Delete("ReturnUri")
	//sessionStore.Save()


	err := a.gServer.HandleAuthorizeRequest(w, r)

	if err != nil {
		fmt.Printf("Verify: %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	}
}

func userAuthorizeHandler(w http.ResponseWriter, r *http.Request) (userID string, err error) {
	store, err := session.Start(r.Context(), w, r)
	if err != nil {
		return
	}
	uid, ok := store.Get("LoggedInUserID")
	if !ok {
		if r.Form == nil {
			r.ParseForm()
		}
		store.Set("ReturnUri", r.Form)
		store.Save()
		w.Header().Set("Location", fmt.Sprintf("/%s/%s", utils.OAUTH2_PREFIX, "login"))
		w.WriteHeader(http.StatusFound)
		return
	}
	userID = uid.(string)
	store.Delete("LoggedInUserID")
	store.Save()

	return
}

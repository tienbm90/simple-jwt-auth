package api

import (
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	oauth2_modesl "github.com/go-oauth2/oauth2/v4/models"
	ginsvr "github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	"github.com/go-session/session"
	"log"

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
	"os"
	"sync"
	"time"
)

type Oauth2API struct {
	gServer *ginsvr.Server
	once    sync.Once
}

// InitServer Initialize the service
func (a *Oauth2API) InitServer(manager oauth2.Manager) *ginsvr.Server {
	a.once.Do(func() {
		a.gServer = ginsvr.NewDefaultServer(manager)
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
		Domain: "http://localhost:9094",
	})

	manager.MapClientStorage(clientStore)
	// Initialize the oauth2 service

	svr := ginsvr.NewServer(ginsvr.NewConfig(), manager)
	svr.SetUserAuthorizationHandler(userAuthorizeHandler)
	api := Oauth2API{
		gServer: svr,
		once:    sync.Once{},
	}

	return api
}

func (a *Oauth2API) Login(c *gin.Context) {
	loginHandler(c.Writer, c.Request)
}

func (a *Oauth2API) Authenicate(c *gin.Context) {
	authHandler(c.Writer, c.Request)
	c.Abort()
}

func (a *Oauth2API) Authorize(c *gin.Context) {
	r := c.Request
	w := c.Writer

	store, err := session.Start(r.Context(), w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	var form url.Values
	if v, ok := store.Get("ReturnUri"); ok {
		form = v.(url.Values)
		log.Printf("METHOD: %s, ReturnUri: %s", c.Request.Method, v)
	}
	r.Form = form
	redirectUri := form.Get("redirect_uri")

	log.Printf("Retrieved stated: %s", c.Request.URL)
	store.Delete("ReturnUri")
	store.Save()

	user_id, err := a.gServer.UserAuthorizationHandler(w, r)

	if err != nil {
		fmt.Printf("Verify: %s", err.Error())
		http.Error(w, err.Error(), http.StatusBadRequest)
	} else {
		if user_id != "" {
			fmt.Printf("Logged in user: %s:%s ưn", user_id, redirectUri)
			http.Redirect(w, r, redirectUri, http.StatusFound)
		}
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
		log.Printf("Set ReturnURI: %s", r.Form)
		store.Save()

		w.Header().Set("Location", "/oauth2/login")
		w.WriteHeader(http.StatusFound)
		return
	}
	state, ok := store.Get("State")
	if ok {
		log.Printf("Staet: %s", state)
	}
	userID = uid.(string)
	store.Delete("LoggedInUserID")
	store.Save()

	return
}

func loginHandler(w http.ResponseWriter, r *http.Request) {
	store, err := session.Start(r.Context(), w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if r.Method == "POST" {
		if r.Form == nil {
			if err := r.ParseForm(); err != nil {
				http.Error(w, err.Error(), http.StatusInternalServerError)
				return
			}
		}
		store.Set("LoggedInUserID", r.Form.Get("username"))
		store.Save()

		w.Header().Set("Location", "/oauth2/auth")
		w.WriteHeader(http.StatusFound)
		return
	}
	outputHTML(w, r, "static/templates/oauth2/login.html")
}

func authHandler(w http.ResponseWriter, r *http.Request) {
	store, err := session.Start(nil, w, r)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if _, ok := store.Get("LoggedInUserID"); !ok {
		log.Printf("User doesnot log in. Loggin!!!!")
		w.Header().Set("Location", "/oauth2/login")

		return
	}

	v, ok := store.Get("ReturnUri")
	if ok {
		log.Printf("Stored ReturnURI: %s", v)
	}
	outputHTML(w, r, "static/templates/oauth2/auth.html")
}

func outputHTML(w http.ResponseWriter, req *http.Request, filename string) {
	file, err := os.Open(filename)
	if err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	defer file.Close()
	fi, _ := file.Stat()
	http.ServeContent(w, req, file.Name(), fi.ModTime(), file)
}

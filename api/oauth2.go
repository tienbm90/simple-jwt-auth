package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/simple-jwt-auth/models"

	//"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/go-oauth2/oauth2/v4/errors"
	"github.com/go-oauth2/oauth2/v4/generates"
	"github.com/go-oauth2/oauth2/v4/manage"
	oauth2_modesl "github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	"github.com/go-session/session"
	"github.com/simple-jwt-auth/ginserver"
	"github.com/simple-jwt-auth/utils"
	"log"
	"net/http"
	"net/url"
	"strings"
	"time"
)

type Oauth2API struct {
	gServer         *server.Server
	oauthClientRepo *models.OauthClientRepository
	userRepository  *models.UserRepository
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
	clientStore.Set("22222", &oauth2_modesl.Client{
		ID:     "22222",
		Secret: "22222222",
		Domain: "http://localhost:8085",
	})
	manager.MapClientStorage(clientStore)
	// Initialize the oauth2 service
	svr := ginserver.InitServer(manager)
	svr.SetPasswordAuthorizationHandler(func(username, password string) (userID string, err error) {
		if username == "test" && password == "test" {
			userID = "test"
		} else {
			userID = "demo"
		}
		return
	})

	svr.SetUserAuthorizationHandler(userAuthorizeHandler)
	svr.SetClientScopeHandler(clientScopeHandler)
	svr.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})
	svr.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	api := Oauth2API{
		gServer: svr,
	}
	return api
}

func InitOauth2API(clientRepository *models.OauthClientRepository, userRepository *models.UserRepository) Oauth2API {
	// init oauth2 server
	manager := manage.NewDefaultManager()
	manager.SetAuthorizeCodeTokenCfg(manage.DefaultAuthorizeCodeTokenCfg)
	// token store
	//manager.MustTokenStorage(store.NewFileTokenStore("data.db"))
	manager.MustTokenStorage(store.NewMemoryTokenStore())
	manager.MapAccessGenerate(generates.NewJWTAccessGenerate("", []byte("00000000"), jwt.SigningMethodHS512))
	// client store
	clientStore := store.NewClientStore()

	clients, err := clientRepository.FindAll()

	if err != nil {
		log.Fatalf("Cannot create client clientRepository")
	}

	for _, v := range clients {
		clientStore.Set(v.ClientID, &oauth2_modesl.Client{
			ID:     v.ClientID,
			Secret: v.ClientSecret,
			Domain: v.Domain,
			UserID: v.UserID,
		})
	}
	//clientStore.Set("22222", &oauth2_modesl.Client{
	//	ID:     "22222",
	//	Secret: "22222222",
	//	Domain: "http://localhost:8085",
	//})
	manager.MapClientStorage(clientStore)
	// Initialize the oauth2 service
	svr := ginserver.InitServer(manager)
	svr.SetPasswordAuthorizationHandler(func(username, password string) (userID string, err error) {
		_, err = userRepository.Validate(models.User{UserName: username, Password: password})
		if err != nil {
			return "", err
		} else {
			return username, nil
		}
	})

	svr.SetUserAuthorizationHandler(userAuthorizeHandler)
	svr.SetClientScopeHandler(clientScopeHandler)
	svr.SetInternalErrorHandler(func(err error) (re *errors.Response) {
		log.Println("Internal Error:", err.Error())
		return
	})
	svr.SetResponseErrorHandler(func(re *errors.Response) {
		log.Println("Response Error:", re.Error.Error())
	})

	api := Oauth2API{
		gServer:         svr,
		oauthClientRepo: clientRepository,
		userRepository:  userRepository,
	}
	return api
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

func (a *Oauth2API) ShowLoginPage(c *gin.Context) {
	w := c.Writer
	store, err := session.Start(c.Request.Context(), w, c.Request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	link := CreateLoginURL(c.Request)
	if c.Request.Form == nil {
		err := c.Request.ParseForm()
		if err != nil {
			log.Fatal("error in parse form")
		} else {

		}
	}
	state := c.Request.Form.Get("state")

	if state != "" {
		store.Set("state", c.Request.Form.Get("state"))
	} else {
		log.Println("Cant find login state!!!")
	}

	store.Save()

	c.HTML(http.StatusOK, "login.tmpl", gin.H{"link": link})
}

func (a *Oauth2API) HandleLogin(c *gin.Context) {
	start := time.Now()
	store, err := session.Start(c.Request.Context(), c.Writer, c.Request)

	if err != nil {
		log.Fatal("Can't init session store")
		return
	}

	if err != nil {
		c.HTML(http.StatusInternalServerError, "login.tmpl", gin.H{
			"ErrorTitle":   "Login Failed",
			"ErrorMessage": "Invalid credentials provided"})
		return
	}

	state, ok := store.Get("state")
	if ok && state != "" {
		log.Printf("LOGIN state: %s \n", state)
	} else {
		log.Printf("LOGIN  state not found: %s \n", state)
		store.Set("state", c.Request.Form.Get("state"))
	}

	if c.Request.Method == "GET" {
		link := CreateLoginURL(c.Request)
		c.HTML(http.StatusOK, "login.tmpl", gin.H{"link": link})
		return
	}



	//store state
	if c.Request.Form == nil {
		err := c.Request.ParseForm()
		if err != nil {
			log.Fatal("error in parse form")
		}
	}

	username := c.PostForm("username")
	password := c.PostForm("password")
	ok, err = a.userRepository.Validate(models.User{UserName: username, Password: password})

	var link string
	if err == nil && ok {
		link = CreateAuthorizeURL(c.Request)
		log.Println("Redirect to url " + link)
		store.Set("LoggedInUserID", username)
		c.Redirect(http.StatusFound, link)

	} else {
		link = CreateLoginURL(c.Request)
		log.Println("Redirect to url " + link)

		c.HTML(http.StatusSeeOther, "login.tmpl", gin.H{
			"ErrorTitle":   "Login Failed",
			"Location":     link,
			"ErrorMessage": "Invalid credentials provided"})
	}

	//link := CreateLoginURL(c.Request)
	//
	//c.HTML(http.StatusFound, "login.tmpl", gin.H{"link": link})

	elapsed := time.Since(start)
	log.Printf("Binomial took %s", elapsed)
	c.Abort()
	//return
}

func (a *Oauth2API) Authenicate(c *gin.Context) {
	store, err := session.Start(c.Request.Context(), c.Writer, c.Request)
	if err != nil {
		http.Error(c.Writer, err.Error(), http.StatusInternalServerError)
		return
	}

	state, ok := store.Get("state")
	if ok {
		log.Printf("Authenticate state: %s \n", state)
	} else {
		log.Printf("Authenticate  state not found: %s \n", state)
	}

	//check if user is logged in

	if _, ok := store.Get("LoggedInUserID"); !ok {
		log.Println("User doesn't log in")

		link := CreateLoginURL(c.Request)

		log.Printf("login link: %s \n", link)
		//c.Writer.Header().Set("Location", fmt.Sprintf("/%s/%s", utils.OAUTH2_PREFIX, "login"))
		c.Writer.Header().Set("Location", link)

		return
	} else {
		log.Println("User is logged in")
	}

	v, ok := store.Get("ReturnUri")
	if ok {
		log.Printf("Stored ReturnURI: %s", v)
	}

	store.Save()
	c.HTML(http.StatusFound, "auth.tmpl", gin.H{"link": "link"})
}

func (a *Oauth2API) Authorize(c *gin.Context) {
	start := time.Now()
	w := c.Writer
	store, err := session.Start(c.Request.Context(), w, c.Request)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if c.Request.Form == nil {
		c.Request.ParseForm()
	}
	state, ok := store.Get("state")
	if ok {
		log.Printf("Pre Authorize state:%s", state)
	} else {
		store.Set("state", c.Request.Form.Get("state"))
		c.SetCookie("state", c.Request.Form.Get("state"), 3600, "", "", false, true)
		log.Printf("Alter Authorize state:%s", c.Request.Form.Get("state"))
	}

	var form url.Values
	form = c.Request.Form
	log.Printf("Form: %s", form)
	//store.Delete("ReturnUri")
	store.Save()

	elapsed := time.Since(start)
	log.Printf("Binomial took %s", elapsed)

	err = a.gServer.HandleAuthorizeRequest(w, c.Request)

	//a.HandleAuthorizeRequest(c)
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
	log.Printf("Request URI: %s", r.RequestURI)
	uid, ok := store.Get("LoggedInUserID")
	if !ok {
		if r.Form == nil {
			r.ParseForm()
		}
		store.Set("ReturnUri", r.Form)
		store.Save()

		link := CreateLoginURL(r)
		//link := "/oauth2/login"
		log.Println("User is not loggin ")

		w.Header().Set("Location", link)
		w.WriteHeader(http.StatusFound)

		return
	} else {
		log.Printf("User: %s \n", uid)
	}
	userID = uid.(string)
	//store.Delete("LoggedInUserID")
	store.Save()
	return
}

func clientScopeHandler(clientID, scope string) (allowed bool, err error) {

	if scope != "all" {
		return false, errors.ErrInvalidGrant
	}
	return true, nil
}

func CreateLoginURL(r *http.Request) string {
	var buf bytes.Buffer
	baseUrl := fmt.Sprintf("/%s/%s", utils.OAUTH2_PREFIX, "login")
	buf.WriteString(baseUrl)
	v := url.Values{
		"response_type": {"code"},
	}

	client_id := r.FormValue("client_id")

	if client_id != "" {
		v.Set("client_id", client_id)
	}

	state := r.FormValue("state")
	if state != "" {
		v.Set("state", state)
	}

	redirect_uri := r.FormValue("redirect_uri")
	if redirect_uri != "" {
		v.Set("redirect_uri", redirect_uri)
	}

	v.Set("scope", "all")

	if strings.Contains(baseUrl, "?") {
		buf.WriteByte('&')
	} else {
		buf.WriteByte('?')
	}
	buf.WriteString(v.Encode())
	return buf.String()
}

func CreateAuthorizeURL(r *http.Request) string {
	var buf bytes.Buffer
	baseUrl := fmt.Sprintf("/%s/%s", utils.OAUTH2_PREFIX, "authorize")
	buf.WriteString(baseUrl)
	v := url.Values{
		"response_type": {"code"},
	}

	client_id := r.FormValue("client_id")
	if client_id != "" {
		v.Set("client_id", client_id)
	}

	state := r.FormValue("state")
	if state != "" {
		v.Set("state", state)
	}
	redirect_uri := r.FormValue("redirect_uri")
	if redirect_uri != "" {
		v.Set("redirect_uri", redirect_uri)
	}

	v.Set("scope", "all")

	if strings.Contains(baseUrl, "?") {
		buf.WriteByte('&')
	} else {
		buf.WriteByte('?')
	}
	buf.WriteString(v.Encode())
	return buf.String()
}

// HandleAuthorizeRequest the authorization request handling
func (a *Oauth2API) HandleAuthorizeRequest(c *gin.Context) {
	err := a.gServer.HandleAuthorizeRequest(c.Writer, c.Request)
	if err != nil {
		c.AbortWithError(http.StatusBadRequest, err)

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

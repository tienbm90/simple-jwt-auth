package api
//
//import (
//	"context"
//	"encoding/json"
//	"fmt"
//	"github.com/gin-gonic/gin"
//	"golang.org/x/oauth2"
//	"golang.org/x/oauth2/clientcredentials"
//	"io"
//	"log"
//	"net/http"
//	"time"
//)
//
//type OauthAPI struct {
//	Config *oauth2.Config
//	globalToken *oauth2.Token
//}
//
//func ProvideOauthAPI(config *oauth2.Config) OauthAPI {
//	return OauthAPI{
//		Config: config,
//	}
//}
//
//const (
//	authServerURL = "http://localhost:9096"
//)
//
//func (g OauthAPI) IndexHandler(c *gin.Context) {
//	r := c.Request
//	w := c.Writer
//	u := g.Config.AuthCodeURL("xyz")
//	log.Println(u)
//	http.Redirect(w, r, u, http.StatusFound)
//	//u := g.Config.AuthCodeURL("xyz")
//	//http.Redirect(w, r, u, http.StatusFound)
//	//c.Redirect(http.StatusFound, u)
//}
//
//// AuthHandler handles authentication of a user and initiates a session.
//func (g OauthAPI) AuthHandler(c *gin.Context) {
//	r := c.Request
//	w := c.Writer
//	r.ParseForm()
//	state := r.Form.Get("state")
//	if state != "xyz" {
//		http.Error(w, "State invalid", http.StatusBadRequest)
//		return
//	}
//	code := r.Form.Get("code")
//	if code == "" {
//		http.Error(w, "Code not found", http.StatusBadRequest)
//		return
//	}
//	token, err := g.Config.Exchange(context.Background(), code)
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//	g.globalToken = token
//
//	e := json.NewEncoder(w)
//	e.SetIndent("", "  ")
//	e.Encode(token)
//}
//
//// LoginHandler handles the login procedure.
//func (g OauthAPI) RefreshHandler(c *gin.Context) {
//	r := c.Request
//	w := c.Writer
//	if g.globalToken == nil {
//		http.Redirect(w, r, "/", http.StatusFound)
//		return
//	}
//
//	g.globalToken.Expiry = time.Now()
//	token, err := g.Config.TokenSource(context.Background(), g.globalToken).Token()
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	g.globalToken = token
//	e := json.NewEncoder(w)
//	e.SetIndent("", "  ")
//	e.Encode(token)
//}
//
//func (g OauthAPI) TryHandler(c *gin.Context) {
//	r := c.Request
//	w := c.Writer
//	if g.globalToken == nil {
//		http.Redirect(w, r, "/", http.StatusFound)
//		return
//	}
//
//	resp, err := http.Get(fmt.Sprintf("%s/test?access_token=%s", authServerURL, g.globalToken.AccessToken))
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//	defer resp.Body.Close()
//
//	io.Copy(w, resp.Body)
//
//}
//
//// FieldHandler is a rudementary handler for logged in users.
//func (g OauthAPI) FwdHandler(c *gin.Context) {
//	r := c.Request
//	w := c.Writer
//	if g.globalToken == nil {
//		http.Redirect(w, r, "/", http.StatusFound)
//		return
//	}
//
//	resp, err := http.Get(fmt.Sprintf("%s/test?access_token=%s", authServerURL, g.globalToken.AccessToken))
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusBadRequest)
//		return
//	}
//	defer resp.Body.Close()
//
//	io.Copy(w, resp.Body)
//}
//
//func (g OauthAPI) ClientHandler(c *gin.Context) {
//	//r:=c.Request
//	w := c.Writer
//	cfg := clientcredentials.Config{
//		ClientID:     g.Config.ClientID,
//		ClientSecret: g.Config.ClientSecret,
//		TokenURL:     g.Config.Endpoint.TokenURL,
//	}
//
//	token, err := cfg.Token(context.Background())
//	if err != nil {
//		http.Error(w, err.Error(), http.StatusInternalServerError)
//		return
//	}
//
//	e := json.NewEncoder(w)
//	e.SetIndent("", "  ")
//	e.Encode(token)
//}

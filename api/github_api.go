package api

import (
	"encoding/json"
	"fmt"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/simple-jwt-auth/middleware"
	"github.com/simple-jwt-auth/models"
	"golang.org/x/oauth2"
	"gorm.io/gorm"
	"io/ioutil"
	"log"
	"net/http"
)

type GithubAPI struct {
	Config   *oauth2.Config
	UserRepo *models.UserRepository
}

func ProvideGithubAPI(config *oauth2.Config, repository *models.UserRepository) GoogleAPI {
	return GoogleAPI{
		Config:   config,
		UserRepo: repository,
	}
}

func (g GithubAPI) IndexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{"link": "/oauth/github/login"})
}

// AuthHandler handles authentication of a user and initiates a session.
func (g GithubAPI) AuthHandler(c *gin.Context) {
	// Handle the exchange code to initiate a transport.
	session := sessions.Default(c)
	retrievedState := session.Get("state")
	queryState := c.Request.URL.Query().Get("state")
	if retrievedState != queryState {
		log.Printf("Invalid session state: retrieved: %s; Param: %s", retrievedState, queryState)
		c.HTML(http.StatusUnauthorized, "error.tmpl", gin.H{"message": "Invalid session state."})
		return
	}

	log.Println(fmt.Sprintf("queryState: %s", queryState))

	code := c.Request.URL.Query().Get("code")
	tok, err := g.Config.Exchange(oauth2.NoContext, code)
	if err != nil {
		log.Println(err)
		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"message": "Login failed. Please try again."})
		return
	}

	client := g.Config.Client(oauth2.NoContext, tok)
	userinfo, err := client.Get("https://api.github.com/user")
	if err != nil {
		log.Println(err)
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	defer userinfo.Body.Close()
	data, _ := ioutil.ReadAll(userinfo.Body)
	githubUser := models.GithubUser{}
	if err = json.Unmarshal(data, &githubUser); err != nil {
		log.Println(err)
		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"message": "Error marshalling response. Please try agian."})
		return
	}

	if githubUser.Email == "" {
		primaryEmail, err := GetPrimaryEmail(client)
		if err != nil {
			githubUser.Email = fmt.Sprintf("%s@gmail.com", githubUser.Login)
		} else {
			githubUser.Email = primaryEmail.Email
			githubUser.EmailVerified = primaryEmail.Verified
		}
	}

	u := models.User{
		Model:          gorm.Model{},
		UserName:       githubUser.Login,
		AuthorizorType: "Github",
		Password:       "",
		Sub:            "",
		Name:           githubUser.Name,
		GivenName:      "",
		FamilyName:     "",
		Profile:        githubUser.Url,
		Picture:        githubUser.AvatarUrl,
		Email:          githubUser.Email,
		EmailVerified:  githubUser.EmailVerified,
		Gender:         "",
	}

	session.Set("user-id", u.Email)
	err = session.Save()
	if err != nil {
		log.Println(err)
		c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"message": "Error while saving session. Please try again."})
		return
	}
	seen := false

	//save user into database if user is not exist
	_, err = g.UserRepo.FindByEmail(u.Email)
	if err != nil {
		_, err = g.UserRepo.Create(u)
		if err != nil {
			log.Println(err)
			c.HTML(http.StatusBadRequest, "error.tmpl", gin.H{"message": "Error while saving user. Please try again."})
			return
		}
		log.Println(fmt.Sprintf("Not found user %s. Create new record", githubUser.Email))
	} else {
		seen = true
	}
	c.HTML(http.StatusOK, "battle.tmpl", gin.H{"email": u.Email, "seen": seen})
}

// LoginHandler handles the login procedure.
func (g GithubAPI) LoginHandler(c *gin.Context) {
	state, err := middleware.RandToken(32)
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{"message": "Error while generating random data."})
		return
	}
	session := sessions.Default(c)
	session.Set("state", state)

	log.Println(fmt.Sprintf("State: %s", state))
	err = session.Save()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "error.tmpl", gin.H{"message": "Error while saving session."})
		return
	}

	// get login url
	link := g.Config.AuthCodeURL(state, oauth2.AccessTypeOnline)
	c.HTML(http.StatusOK, "auth.tmpl", gin.H{"link": link})
}

func (g GithubAPI) TestHandler(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user-id")
	c.HTML(http.StatusOK, "field.tmpl", gin.H{"user": userID, "link": "/oauth/field"})

}

// FieldHandler is a rudementary handler for logged in users.
func (g GithubAPI) FieldHandler(c *gin.Context) {
	session := sessions.Default(c)
	userID := session.Get("user-id")
	c.HTML(http.StatusOK, "field.tmpl", gin.H{"user": userID})
}

func (g GithubAPI) ApiHandler(c *gin.Context) {
	c.JSON(http.StatusOK, "Ok test")
}

func GetPrimaryEmail(client *http.Client) (models.GithubEmail, error) {
	userinfo, err := client.Get("https://api.github.com/user")
	if err != nil {
		log.Println(err)
		return models.GithubEmail{}, err
	}
	defer userinfo.Body.Close()
	data, _ := ioutil.ReadAll(userinfo.Body)
	githubEmails := models.GithubEmails{}
	if err = json.Unmarshal(data, &githubEmails); err != nil {
		log.Println(err)
		return models.GithubEmail{}, err
	}
	for _, ema := range githubEmails {
		if ema.Primary == true {
			return ema, nil
		}
	}
	return models.GithubEmail{}, err
}

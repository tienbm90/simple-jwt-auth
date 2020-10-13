package api

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/simple-jwt-auth/auth"
	"github.com/simple-jwt-auth/models"
	"log"
	"net/http"
	"os"
	"strconv"
)

var tokenManager = auth.JwtTokenManager{}

type JwtApi struct {
	UserRepo *models.UserRepository
}

func CreateJwtApi(repo *models.UserRepository) JwtApi {
	return JwtApi{UserRepo: repo}
}

func (api JwtApi) JwtLogin(c *gin.Context) {
	var u models.User
	if err := c.ShouldBindJSON(&u); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	user, err := api.UserRepo.Validate(u)

	if err != nil {
		c.JSON(http.StatusUnauthorized, "Please provide valid login details")
		return
	}
	log.Println(user)

	ts, err := tokenManager.CreateToken(fmt.Sprintf("%s",user.ID), user.UserName)
	if err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}

	// save token to redis
	//saveErr := servers.HttpServer.RD.CreateAuth(user.ID, ts)
	//if saveErr != nil {
	//	c.JSON(http.StatusUnprocessableEntity, saveErr.Error())
	//}
	tokens := map[string]string{
		"access_token":  ts.AccessToken,
		"refresh_token": ts.RefreshToken,
	}
	c.JSON(http.StatusOK, tokens)
}

func (api JwtApi) JwtLogout(c *gin.Context) {
	//If metadata is passed and the tokens valid, delete them from the redis store
	metadata, _ := tokenManager.ExtractTokenMetadata(c.Request)
	if metadata != nil {
		//deleteErr := servers.HttpServer.RD.DeleteTokens(metadata)
		//if deleteErr != nil {
		//	c.JSON(http.StatusBadRequest, deleteErr.Error())
		//	return
		//}
	}
	c.JSON(http.StatusOK, "Successfully logged out")
}

func (api JwtApi) JwtRefresh(c *gin.Context) {
	mapToken := map[string]string{}
	if err := c.ShouldBindJSON(&mapToken); err != nil {
		c.JSON(http.StatusUnprocessableEntity, err.Error())
		return
	}
	refreshToken := mapToken["refresh_token"]

	//verify the token
	token, err := jwt.Parse(refreshToken, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return []byte(os.Getenv("REFRESH_SECRET")), nil
	})
	//if there is an error, the token must have expired
	if err != nil {
		c.JSON(http.StatusUnauthorized, "JwtRefresh token expired")
		return
	}
	//is token valid?
	if _, ok := token.Claims.(jwt.Claims); !ok && !token.Valid {
		c.JSON(http.StatusUnauthorized, err)
		return
	}
	//Since token is valid, get the uuid:
	claims, ok := token.Claims.(jwt.MapClaims) //the token claims should conform to MapClaims
	if ok && token.Valid {
		//refreshUuid, ok := claims["refresh_uuid"].(string) //convert the interface to string
		//if !ok {
		//	c.JSON(http.StatusUnprocessableEntity, err)
		//	return
		//}
		userId, roleOk := claims["user_id"].(string)
		if roleOk == false {
			c.JSON(http.StatusUnprocessableEntity, "unauthorized")
			return
		}
		//Delete the previous JwtRefresh Token
		//delErr := servers.HttpServer.RD.DeleteRefresh(refreshUuid)
		//if delErr != nil { //if any goes wrong
		//	c.JSON(http.StatusUnauthorized, "unauthorized")
		//	return
		//}
		//Create new pairs of refresh and access tokens

		userID, err := strconv.Atoi(userId)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, "userId invalid")
			return
		}
		user, err := api.UserRepo.FindByID(userID)
		if err != nil {
			c.JSON(http.StatusUnprocessableEntity, "Subject's not found ")
		}

		ts, createErr := tokenManager.CreateToken(userId, user.UserName)
		if createErr != nil {
			c.JSON(http.StatusForbidden, createErr.Error())
			return
		}
		//save the tokens metadata to redis
		//saveErr := servers.HttpServer.RD.CreateAuth(userId, ts)
		//if saveErr != nil {
		//	c.JSON(http.StatusForbidden, saveErr.Error())
		//	return
		//}
		tokens := map[string]string{
			"access_token":  ts.AccessToken,
			"refresh_token": ts.RefreshToken,
		}
		c.JSON(http.StatusCreated, tokens)
	} else {
		c.JSON(http.StatusUnauthorized, "refresh expired")
	}
}

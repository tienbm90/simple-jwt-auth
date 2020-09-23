package middleware

import (
	"errors"
	"fmt"
	"github.com/casbin/casbin"
	"github.com/casbin/casbin/persist"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/gophercon-jwt-repo/auth"
	"log"
	"net/http"
)

func TokenAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		err := auth.TokenValid(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "unauthorized")
			c.Abort()
			return
		}
		c.Next()
	}
}

func ExtractTokenMetadata(r *http.Request) (*auth.AccessDetails, error) {
	token, err := auth.VerifyToken(r)
	if err != nil {
		return nil, err
	}
	acc, err := extract(token)
	if err != nil {
		return nil, err
	}
	return acc, nil
}

func extract(token *jwt.Token) (*auth.AccessDetails, error) {

	claims, ok := token.Claims.(jwt.MapClaims)
	if ok && token.Valid {
		accessUuid, ok := claims["access_uuid"].(string)
		userId, userOk := claims["user_id"].(string)
		userName, userNameOk := claims["user_name"].(string)
		if ok == false || userOk == false || userNameOk == false {
			return nil, errors.New("unauthorized")
		} else {
			return &auth.AccessDetails{
				TokenUuid: accessUuid,
				UserId:    userId,
				UserName:  userName,
			}, nil
		}
	}
	return nil, errors.New("something went wrong")
}

// Authorize determines if current subject has been authorized to take an action on an object.
func Authorize(obj string, act string, adapter persist.Adapter) gin.HandlerFunc {
	return func(c *gin.Context) {
		//val, existed := c.Get("current_subject")
		//if !existed {
		//	c.AbortWithStatusJSON(401, "user hasn't logged in yet")
		//	return
		//}

		err := auth.TokenValid(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "user hasn't logged in yet")
			c.Abort()
			return
		}
		metadata, err := ExtractTokenMetadata(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "unauthorized")
			return
		}
		userId, err := h.rd.FetchAuth(metadata.TokenUuid)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "unauthorized")
			return
		}

		// casbin enforces policy
		ok, err := enforce("admin", obj, act, adapter)
		//ok, err := enforce(val.(string), obj, act, adapter)
		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(500, "error occurred when authorizing user")
			return
		}
		if !ok {
			c.AbortWithStatusJSON(403, "forbidden")
			return
		}
		c.Next()
	}
}

func enforce(sub string, obj string, act string, adapter persist.Adapter) (bool, error) {
	enforcer := casbin.NewEnforcer("config/rbac_model.conf", adapter)
	//if err != nil {
	//	return false, fmt.Errorf("failed to create casbin enforcer: %w", err)
	//}
	// Load policies from DB dynamically
	err := enforcer.LoadPolicy()
	if err != nil {
		return false, fmt.Errorf("failed to load policy from DB: %w", err)
	}
	ok := enforcer.Enforce(sub, obj, act)
	return ok, err
}



func  FetchAuth(tokenUuid string) (string, error) {
	userid, err := .Get(tokenUuid).Result()
	if err != nil {
		return "", err
	}
	return userid, nil
}
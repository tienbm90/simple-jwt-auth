package middleware

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/simple-jwt-auth/auth"
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

// Authorize determines if current subject has been authorized to take an action on an object.
func Authorize(obj string, act string, enforcer *casbin.Enforcer) gin.HandlerFunc {
	return func(c *gin.Context) {
		err := auth.TokenValid(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "user hasn't logged in yet")
			c.Abort()
			return
		}
		metadata, err := auth.ExtractTokenMetadata(c.Request)
		if err != nil {
			c.JSON(http.StatusUnauthorized, "unauthorized")
			return
		}
		// casbin enforces policy
		log.Printf("Meta: %s:%s:%s", metadata.UserName, obj, act)
		ok, err := enforce(metadata.UserName, obj, act, enforcer)
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

func enforce(sub string, obj string, act string, enforcer *casbin.Enforcer) (bool, error) {
	ok, err := enforcer.Enforce(sub, obj, act)
	return ok, err
}

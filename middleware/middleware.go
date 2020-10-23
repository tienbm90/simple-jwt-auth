package middleware

import (
	"crypto/rand"
	"encoding/base64"
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/simple-jwt-auth/middleware/auth"
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

// AuthorizeJwtToken determines if current subject has been authorized to take an action on an object.
func AuthorizeJwtToken(obj string, act string, enforcer *casbin.Enforcer) gin.HandlerFunc {
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
		log.Println("Meta: %s:%s:%s", metadata.UserName, obj, act)
		ok, err := enforce(metadata.UserName, obj, act, enforcer)
		//ok, err := enforce(val.(string), obj, act, adapter)

		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(500, "error occurred when authorizing user")
			return
		}
		if !ok {
			c.AbortWithStatusJSON(403, "Permission Invalid")
			return
		}
		c.Next()
	}
}

func enforce(sub string, obj string, act string, enforcer *casbin.Enforcer) (bool, error) {
	ok, err := enforcer.Enforce(sub, obj, act)
	return ok, err
}


// AuthorizeRequest is used to authorize a request for a certain end-point group.
func AuthorizeOpenIdRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		session := sessions.Default(c)
		v := session.Get("user-id")
		if v == nil {
			c.HTML(http.StatusUnauthorized, "error.tmpl", gin.H{"message": "Please login."})
			c.Abort()
		}
		c.Next()
	}
}


// RandToken generates a random @l length token.
func RandToken(l int) (string, error) {
	b := make([]byte, l)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return base64.StdEncoding.EncodeToString(b), nil
}
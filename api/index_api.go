package api

import (
	"github.com/gin-gonic/gin"
	"net/http"
)


// IndexHandler handles the location /.
func IndexHandler(c *gin.Context) {
	c.HTML(http.StatusOK, "index.tmpl", gin.H{})
}

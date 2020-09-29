package api

import (
	"github.com/gin-gonic/gin"
	"github.com/simple-jwt-auth/auth"
	"github.com/simple-jwt-auth/models"
	"net/http"
)

func CreateTodo(c *gin.Context) {
	var td models.Todo
	if err := c.ShouldBindJSON(&td); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "invalid json")
		return
	}
	metadata, err := auth.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	td.UserID = metadata.UserId

	//you can proceed to save the  to a database

	c.JSON(http.StatusCreated, td)
}

func GetTodo(c *gin.Context) {
	metadata, err := auth.ExtractTokenMetadata(c.Request)
	if err != nil {
		c.JSON(http.StatusUnauthorized, "unauthorized")
		return
	}

	userId := metadata.UserId
	c.JSON(http.StatusOK, models.Todo{
		UserID: userId,
		Title:  "Return todo",
		Body:   "Return todo for testing",
	})
}

package api

import (
	"fmt"
	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/simple-jwt-auth/models"
	"log"
	"net/http"
	"strconv"
)

type UserAPI struct {
	UserService UserService
}

func ProvideUserAPI(s UserService) UserAPI {
	return UserAPI{UserService: s}
}

func (a *UserAPI) FindAll(c *gin.Context) {
	users, err := a.UserService.FindAll()
	if err != nil {
		log.Fatalln(err)
		c.Status(http.StatusBadRequest)
		return
	}

	c.JSON(http.StatusOK, gin.H{"users": ToUserDTOs(users)})
}

func (a *UserAPI) UserInfo(c *gin.Context) {
	session := sessions.Default(c)
	v := session.Get("user-id")

	if v == nil {
		c.HTML(http.StatusUnauthorized, "error.tmpl", gin.H{"message": "Please login."})
		c.Abort()
	}

	id, _ := strconv.Atoi(fmt.Sprint(v))

	user, err := a.UserService.FindByID(int(id))
	if err != nil {
		log.Fatalln(err)
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": ToUserDTO(user)})
}

func (a *UserAPI) FindByID(c *gin.Context) {

	id, _ := strconv.Atoi(c.Param("id"))

	user, err := a.UserService.FindByID(int(id))
	if err != nil {
		log.Fatalln(err)
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": ToUserDTO(user)})
}

func (a *UserAPI) Create(c *gin.Context) {
	var userDTO UserDTO
	err := c.BindJSON(&userDTO)
	if err != nil {
		log.Fatalln(err)
		c.Status(http.StatusBadRequest)
		return
	}
	createdUser, err := a.UserService.Create(ToUser(userDTO))
	if err != nil {
		log.Fatalln(err)
		c.Status(http.StatusBadRequest)
		return
	}
	c.JSON(http.StatusOK, gin.H{"user": ToUserDTO(createdUser)})
}

func (api *UserAPI) Update(c *gin.Context) {
	var userDTO UserDTO
	err := c.BindJSON(&userDTO)
	if err != nil {
		log.Fatalln(err)
		c.Status(http.StatusBadRequest)
		return
	}

	id, _ := strconv.Atoi(c.Param("id"))
	user, err := api.UserService.FindByID(int(id))
	if err != nil {
		log.Fatalln(err)
		c.Status(http.StatusBadRequest)
		return
	}
	if user == (models.User{}) {
		c.Status(http.StatusBadRequest)
		return
	}

	user.UserName = userDTO.Username
	user.Email = userDTO.Email
	user.Password = userDTO.Password
	api.UserService.Update(user)

	c.Status(http.StatusOK)
}

func (api *UserAPI) Delete(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	user, err := api.UserService.FindByID(id)
	if err != nil {
		log.Fatalln(err)
		c.Status(http.StatusBadRequest)
		return
	}
	if user == (models.User{}) {
		c.Status(http.StatusBadRequest)
		return
	}
	api.UserService.Delete(user)
	c.Status(http.StatusOK)
}

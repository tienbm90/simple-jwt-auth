package api

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/simple-jwt-auth/models"
	"net/http"
)

type CasbinService struct {
	Enforcer *casbin.Enforcer
}

func (service *CasbinService) CreatePolicy(c *gin.Context) {
	var p models.Policy
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	ok, err := service.Enforcer.AddPolicy(p.User, p.Path, p.Method)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}
	c.JSON(http.StatusOK, ok)
}

func (service *CasbinService) CreateGroupPolicy(c *gin.Context) {
	var p models.GroupPolicy
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	ok, err := service.Enforcer.AddGroupingPolicy(p.Member, p.Group``)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}
	c.JSON(http.StatusOK, ok)
}

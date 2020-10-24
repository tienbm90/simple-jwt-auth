package api

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/simple-jwt-auth/models"
	"log"
	"net/http"
)

type CasbinService struct {
	Enforcer *casbin.Enforcer
}

func NewCasbinService(enforcer *casbin.Enforcer) *CasbinService {
	return &CasbinService{Enforcer: enforcer}
}

func (service *CasbinService) ListPolicy(c *gin.Context) {
	//reload policies list
	service.Enforcer.LoadPolicy()
	policy := service.Enforcer.GetPolicy()
	c.JSON(http.StatusOK, policy)
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

func (service *CasbinService) DeletePolicy(c *gin.Context) {
	var p models.Policy
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	_, err := service.Enforcer.RemovePolicy(p.User, p.Path, p.Method)
	if err != nil {
		log.Println(err.Error())
		c.JSON(http.StatusBadRequest, err)
	}

	service.Enforcer.LoadPolicy()
	c.JSON(http.StatusOK, "Policies delete successful")
}

func (service *CasbinService) CreateGroupPolicy(c *gin.Context) {
	var p models.GroupPolicy
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	_, err := service.Enforcer.AddGroupingPolicy(p.Member, p.Group)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
		return
	}

	c.JSON(http.StatusOK, p)
}

func (service *CasbinService) ListGroupPolicies(c *gin.Context) {
	groups := service.Enforcer.GetGroupingPolicy()
	c.JSON(http.StatusOK, groups)
}

func (service *CasbinService) DeleteGroupPolicy(c *gin.Context) {
	var p models.GroupPolicy
	if err := c.ShouldBindJSON(&p); err != nil {
		c.JSON(http.StatusUnprocessableEntity, "Invalid json provided")
		return
	}

	ok, err := service.Enforcer.RemoveGroupingPolicy(p.Member, p.Group)
	if err != nil {
		c.JSON(http.StatusBadRequest, err)
	}
	c.JSON(http.StatusOK, ok)
}

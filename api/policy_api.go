package api

import (
	"github.com/casbin/casbin/v2"
	"github.com/gin-gonic/gin"
	"github.com/simple-jwt-auth/middleware/auth"
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

type CasbinAPI struct {
	PoliciesRepo *models.PolicyRepository
}


func CreateCasbinApi(repo *models.PolicyRepository) CasbinAPI {
	return CasbinAPI{PoliciesRepo: repo}
}
func (a *CasbinAPI)ListCasbinRule(c *gin.Context) {
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
func (a *CasbinAPI)CreateCasbinRule(c *gin.Context) {
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

func (a *CasbinAPI)DeleteCasbinRule(c *gin.Context) {
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

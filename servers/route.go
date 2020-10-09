package servers

import (
	"github.com/gin-gonic/gin"
	"github.com/simple-jwt-auth/api"
	"github.com/simple-jwt-auth/middleware"
)

func (s *Server) InitializeRoutes() {

	casbinService := api.NewCasbinService(s.Enforcer)
	s.Router.POST("/login", api.Login)

	authorized := s.Router.Group("/")
	authorized.Use(gin.Logger())
	authorized.Use(gin.Recovery())
	authorized.Use(middleware.TokenAuthMiddleware())
	{
		authorized.POST("/auth/policy",middleware.Authorize("/auth/policy", "POST", s.Enforcer), casbinService.CreatePolicy)
		authorized.GET("/auth/policy", middleware.Authorize("/auth/policy", "GET", s.Enforcer), casbinService.ListPolicy)
		authorized.POST("/auth/grouppolicy",middleware.Authorize("/auth/grouppolicy", "POST", s.Enforcer), casbinService.CreateGroupPolicy)
		authorized.GET("/auth/grouppolicy", middleware.Authorize("/auth/grouppolicy", "GET", s.Enforcer), casbinService.ListGroupPolicies)
		authorized.POST("/api/todo", middleware.Authorize("resource", "write", s.Enforcer), api.CreateTodo)
		authorized.GET("/api/todo", middleware.Authorize("resource", "read", s.Enforcer), api.GetTodo)
		authorized.POST("/logout", api.Logout)
		authorized.POST("/refresh", api.Refresh)
	}

}

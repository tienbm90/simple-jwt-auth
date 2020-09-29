package servers

import (
	"github.com/gin-gonic/gin"
	"github.com/simple-jwt-auth/api"
	"github.com/simple-jwt-auth/middleware"

)

func (s *Server) InitializeRoutes() {

	s.Router.POST("/login", api.Login)

	authorized := s.Router.Group("/")
	authorized.Use(gin.Logger())
	authorized.Use(gin.Recovery())
	authorized.Use(middleware.TokenAuthMiddleware())
	{
		authorized.POST("/api/todo",  middleware.Authorize("resource", "write", s.FileAdapter), api.CreateTodo)
		authorized.GET("/api/todo", middleware.Authorize("resource", "read", s.FileAdapter), api.GetTodo)
		authorized.POST("/logout", api.Logout)
		authorized.POST("/refresh", api.Refresh)
	}

}



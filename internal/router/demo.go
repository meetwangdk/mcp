package router

import (
	"github.com/gin-gonic/gin"
)

func InitUserRouter(
	deps RouterDeps,
	r *gin.RouterGroup,
) {
	// No route group has permission
	noAuthRouter := r.Group("/")
	{
		noAuthRouter.POST("/register", deps.DemoHandler.Register)
		noAuthRouter.POST("/login", deps.DemoHandler.Login)
	}
}

package http

import (
	"github.com/gin-gonic/gin"
)

func SetupRouter(handler *UserHandler) *gin.Engine {
	InitOAuth()

	r := gin.Default()

	r.GET("/auth/login", LoginHandler)
	r.GET("/auth/callback", CallbackHandler)

	protected := r.Group("/")
	protected.Use(AuthMiddleware())
	{
		protected.POST("/users", handler.CreateUser)
		protected.GET("/users/:id", handler.GetUser)
		protected.PUT("/users/:id", handler.UpdateUser)
		protected.DELETE("/users/:id", handler.DeleteUser)
	}

	return r
}

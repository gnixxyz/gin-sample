package auth

import "github.com/gin-gonic/gin"

func Register(r *gin.RouterGroup)  {
	auth := r.Group("/auth")
	{
		auth.POST("/register", register)
		auth.POST("/login", login)
		auth.GET("/check", check)
	}
}
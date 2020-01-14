package v1

import (
	"gin-sample/api/v1/auth"
	"gin-sample/api/v1/posts"
	"github.com/gin-gonic/gin"
	"net/http"
)

func ping(c *gin.Context)  {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}

// Register router to the gin engine
func Register(r *gin.RouterGroup)  {
	v1 := r.Group("/v1")
	{
		v1.GET("/ping", ping)

		auth.Register(v1)
		posts.Register(v1)
	}
}
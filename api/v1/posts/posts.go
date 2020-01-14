package posts

import "github.com/gin-gonic/gin"

func Register(r *gin.RouterGroup)  {
	posts := r.Group("/posts")
	{
		posts.POST("/", create)
		posts.GET("/", list)
		posts.GET("/:id", read)
		posts.DELETE("/:id", remove)
		posts.PATCH("/:id", update)
	}
}
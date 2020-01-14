package views

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func index(c *gin.Context)  {
	c.HTML(http.StatusOK, "index.html", gin.H{})
}

func Routes(r *gin.Engine)  {
	r.LoadHTMLGlob("templates/*")

	views := r.Group("")
	{
		views.GET("/", index)
	}
}
package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func EveryAuth(c *gin.Context)  {
	_, exists := c.Get("user")

	if !exists {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
}
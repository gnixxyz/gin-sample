package api

import (
	v1 "gin-sample/api/v1"
	"github.com/gin-gonic/gin"
)

func Routes(r *gin.Engine)  {
	api := r.Group("/api")
	{
		v1.Register(api)
	}
}
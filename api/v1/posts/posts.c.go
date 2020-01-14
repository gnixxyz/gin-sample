package posts

import (
	"gin-sample/models"
	"gin-sample/pkg/common"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"net/http"
)

type Post = models.Post
type User = models.User
type JSON = common.JSON

func create(c *gin.Context)  {
	db := c.MustGet("db").(*gorm.DB)

	type RequestBody struct {
		Title string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}
	var requestBody RequestBody

	if err := c.BindJSON(&requestBody); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	user := c.MustGet("user").(User)
	post := Post{Title: requestBody.Title, Content:requestBody.Content, User:user}
	db.NewRecord(post)
	db.Create(&post)

	c.JSON(http.StatusOK, post.Serializes())
}

func list(c *gin.Context)  {
	db := c.MustGet("db").(*gorm.DB)

	cursor := c.Query("cursor")
	recent := c.Query("recent")

	var posts []Post

	if cursor == "" {
		if err := db.Preload("User").Limit(10).Order("id desc").Find(&posts).Error; err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	} else {
		condition := "id < ?"
		if recent == "1" {
			condition = "id > ?"
		}
		if err := db.Preload("User").Limit(10).Order("id desc").Where(condition, cursor).Find(&posts).Error; err != nil {
			c.AbortWithStatus(http.StatusInternalServerError)
			return
		}
	}

	size := len(posts)
	serialized := make([]JSON, size, size)

	for i := 0; i < size; i++ {
		serialized[i] = posts[i].Serializes()
	}

	c.JSON(http.StatusOK, serialized)
}

func read(c *gin.Context)  {
	db := c.MustGet("db").(*gorm.DB)

	id := c.Param("id")
	var post Post

	// auto preloads the related model
	// http://gorm.io/docs/preload.html#Auto-Preloading
	if err := db.Set("gorm:auto_preload", true).Where("id = ?", id).First(&post).Error; err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, post.Serializes())
}

func remove(c *gin.Context)  {
	db := c.MustGet("db").(*gorm.DB)

	id := c.Param("id")
	user := c.MustGet("user").(User)

	var post Post
	if err := db.Where("id = ?", id).First(&post).Error; err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	if post.UserID != user.ID {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	db.Delete(&post)
	c.Status(http.StatusNoContent)
}

func update(c *gin.Context)  {
	db := c.MustGet("db").(*gorm.DB)

	id := c.Param("id")

	user := c.MustGet("user").(User)

	type RequestBody struct {
		Title string `json:"title" binding:"required"`
		Content string `json:"content" binding:"required"`
	}
	var requestBody RequestBody

	if err := c.BindJSON(&requestBody); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var post Post
	if err := db.Preload("User").Where("id = ?", id).First(&post).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if post.UserID != user.ID {
		c.AbortWithStatus(http.StatusForbidden)
		return
	}

	post.Title = requestBody.Title
	post.Content = requestBody.Content

	db.Save(&post)
	c.JSON(http.StatusOK, post.Serializes())
}
package models

import (
	"gin-sample/pkg/common"
	"github.com/jinzhu/gorm"
)

type Post struct {
	gorm.Model

	Title string
	Content string `sql:"type:text"`
	User User `gorm:foreignkey:UserID`
	UserID uint
}

func (p Post) Serializes() common.JSON {
	return common.JSON{
		"id": p.ID,
		"title": p.Title,
		"content": p.Content,
		"created_at": p.CreatedAt,
	}
}
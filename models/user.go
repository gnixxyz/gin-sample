package models

import (
	"gin-sample/pkg/common"
	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model

	Username string
	Password string
}

func (u *User) Serializes() common.JSON {
	return common.JSON{
		"id": u.ID,
		"username": u.Username,
		"created_at": u.CreatedAt,
	}
}

func (u *User) Read(m common.JSON) {
	u.ID = uint(m["id"].(float64))
	u.Username = m["username"].(string)
}
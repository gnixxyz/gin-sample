package models

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"os"
)

func Init() (*gorm.DB, error) {
	dbConfig := os.Getenv("DB_CONFIG")
	fmt.Println("DB_CONFIG: ", dbConfig)
	db, err := gorm.Open("mysql", dbConfig)
	if err != nil {
		panic(err)
	}

	db.LogMode(true)

	fmt.Println("Connected to database")

	Migrate(db)

	return db, err
}

func Migrate(db *gorm.DB)  {
	db.AutoMigrate(&User{}, &Post{})
	db.Model(&Post{}).AddForeignKey("user_id", "users(id)", "CASCADE", "CASCADE")
	fmt.Println("Auto Migration has been processed")
}

func Inject(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Set("db", db)
		c.Next()
	}
}
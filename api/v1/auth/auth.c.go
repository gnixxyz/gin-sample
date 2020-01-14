package auth

import (
	"fmt"
	"gin-sample/models"
	"gin-sample/pkg/common"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"golang.org/x/crypto/bcrypt"
	"io/ioutil"
	"net/http"
	"os"
	"time"
)

type User = models.User

func hash(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
	return string(bytes), err
}

func checkPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func generateToken(data common.JSON) (string, error) {
	// token is valid for 7days
	date := time.Now().Add(time.Hour * 24 * 7)

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodES256, jwt.MapClaims{
		"user": data,
		"exp": date.Unix(),
	})

	// get path from root dir
	pwd, _ := os.Getwd()
	keyPath := pwd + "/jwtsecret.key"
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		fmt.Println("fail read jwt secret key file")
		return "", err
	}

	signedToken, err := tokenClaims.SignedString(key)
	fmt.Println("token: ", signedToken)
	return signedToken, err
}

func register(c *gin.Context)  {
	db := c.MustGet("db").(*gorm.DB)

	type RequestBody struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	var requestBody RequestBody
	if err := c.BindJSON(&requestBody); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var exists User
	if err := db.Where("username = ?", requestBody.Username).First(&exists).Error; err == nil {
		c.AbortWithStatus(http.StatusConflict)
		return
	}

	hashPassword, err := hash(requestBody.Password)
	if err != nil {
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	// create
	user := User{
		Username:requestBody.Username,
		Password:hashPassword,
	}

	db.NewRecord(user)
	db.Create(&user)

	serialized := user.Serializes()
	token, _ := generateToken(serialized)

	c.SetCookie("token", token, 60*60*24*7, "/", "", false, true)

	c.JSON(http.StatusOK, common.JSON{
		"user": serialized,
		"token": token,
	})
}

func login(c *gin.Context)  {
	db := c.MustGet("db").(*gorm.DB)

	type RequestBody struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	var requestBody RequestBody
	if err := c.BindJSON(&requestBody); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}

	var user User
	if err := db.Where("username = ?", requestBody.Username).First(&user).Error; err != nil {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}

	if !checkPassword(requestBody.Password, user.Password) {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	serialized := user.Serializes()
	token, _ := generateToken(serialized)

	c.SetCookie("token", token, 60*60*24*7, "/", "", false, true)

	c.JSON(http.StatusOK, common.JSON{
		"user": serialized,
		"token": token,
	})
}

func check(c *gin.Context)  {
	userRaw, ok := c.Get("user")
	if !ok {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}

	user := userRaw.(User)

	tokenExpire := int64(c.MustGet("token_expire").(float64))
	now := time.Now().Unix()
	diff := tokenExpire - now

	fmt.Println(diff)

	if diff < 60*60*24*3 {
		// renew token
		token, _ := generateToken(user.Serializes())
		c.SetCookie("token", token, 60*60*24*7, "/", "", false, true)
		c.JSON(http.StatusOK, common.JSON{
			"token": token,
			"user": user.Serializes(),
		})
		return
	}

	c.JSON(http.StatusOK, common.JSON{
		"token": nil,
		"user": user.Serializes(),
	})
}




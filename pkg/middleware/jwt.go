package middleware

import (
	"errors"
	"fmt"
	"gin-sample/models"
	"gin-sample/pkg/common"
	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"io/ioutil"
	"os"
	"strings"
)

var secretKey []byte

func init()  {
	pwd, _ := os.Getwd()
	keyPath := pwd + "/jwtsecret.key"

	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		panic("failed to load secret key file")
	}

	secretKey = key
}

func validateToken(value string) (common.JSON, error) {
	token, err := jwt.Parse(value, func(token *jwt.Token) (i interface{}, err error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}

		return secretKey, nil
	})

	if err != nil {
		return common.JSON{}, err
	}

	if !token.Valid {
		return common.JSON{}, errors.New("Invalid token")
	}

	return token.Claims.(jwt.MapClaims), nil
}

func JWTMiddleware() gin.HandlerFunc  {
	return func(c *gin.Context) {
		token, err := c.Cookie("token")

		if err != nil {
			authorization := c.Request.Header.Get("Authorization")
			if authorization == "" {
				c.Next()
				return
			}

			tokenArray := strings.Split(authorization, "Bearer ")
			// invalid token
			if len(tokenArray) < 1 {
				c.Next()
				return
			}

			token = tokenArray[1]
		}

		tokenData, err := validateToken(token)
		if err != nil {
			c.Next()
			return
		}

		var user models.User

		user.Read(tokenData["user"].(common.JSON))

		c.Set("user", user)
		c.Set("token_expire", tokenData["exp"])
		c.Next()
	}
}
package utils

import (
	userModels "FileStorage/database/models/user"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
	"log"
	"net/http"
	"time"
)

func GenerateJWT(userID uint, secret []byte) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["user_id"] = userID
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	signedToken, err := token.SignedString(secret)
	if err != nil {
		return "", err
	}

	return signedToken, nil
}

type Claims struct {
	UserID uint `json:"user_id"`
	jwt.StandardClaims
}

func ParseJWT(tokenString string, secret []byte) (*Claims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return secret, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*Claims)
	if !ok || !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	return claims, nil
}

func (claims *Claims) IsExpired() bool {
	return claims.ExpiresAt < time.Now().Unix()
}

func GetCurrentUser(db *gorm.DB, c *gin.Context) (userModels.User, int, error) {
	op := "utils.GetCurrentUser"

	userID, exists := c.Get("user_id")
	if !exists {
		log.Printf("%s: unauthorized", op)
		return userModels.User{}, http.StatusUnauthorized, fmt.Errorf("unauthorized")
	}

	var user userModels.User

	if err := db.Where("id = ?", userID).First(&user).Error; err != nil {
		log.Printf("%s: user not found", op)
		return userModels.User{}, http.StatusNotFound, fmt.Errorf("user not found")
	}

	return user, http.StatusOK, nil
}

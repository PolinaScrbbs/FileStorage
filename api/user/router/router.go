package router

import (
	"FileStorage/api/user/schemes"
	"FileStorage/database/models/user"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func Register(db *gorm.DB) gin.HandlerFunc {
	op := "router.Register"

	return func(c *gin.Context) {
		var reg schemes.RegisterRequest
		if err := c.ShouldBindJSON(&reg); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid input data",
			})
			log.Printf("%s: invalid input data: %v", op, err)
			return
		}

		var existingUser user.User
		if err := db.Where("username = ?", reg.Username).First(&existingUser).Error; err == nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Username already exists",
			})
			log.Printf("%s: username %s already exists", op, reg.Username)
			return
		}

		hash, err := bcrypt.GenerateFromPassword([]byte(reg.Password), bcrypt.DefaultCost)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to hash password",
			})
			log.Printf("%s: error hashing password: %v", op, err)
			return
		}

		new_user := user.User{
			Username:     reg.Username,
			PasswordHash: string(hash),
		}

		if err := db.Create(&new_user).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to register user",
			})
			log.Printf("%s: error creating user: %v", err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "User " + reg.Username + " registered successfully",
		})
		log.Printf("%s: user %s registered successfully", op, reg.Username)
	}
}

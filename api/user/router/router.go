package router

import (
	userSchemes "FileStorage/api/user/schemes"
	"FileStorage/api/user/utils"
	"FileStorage/config"
	tokenModels "FileStorage/database/models/token"
	userModels "FileStorage/database/models/user"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
	"log"
	"net/http"
)

func Register(db *gorm.DB) gin.HandlerFunc {
	op := "router.Register"

	return func(c *gin.Context) {
		var reg userSchemes.UserRequest
		if err := c.ShouldBindJSON(&reg); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid input data",
			})
			log.Printf("%s: invalid input data: %v", op, err)
			return
		}

		var existingUser userModels.User
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

		new_user := userModels.User{
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

func Login(db *gorm.DB, conf *config.Config) gin.HandlerFunc {
	op := "router.Login"

	return func(c *gin.Context) {
		var login userSchemes.UserRequest
		if err := c.ShouldBindJSON(&login); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "Invalid input data",
			})
			log.Printf("%s: invalid input data: %v", op, err)
			return
		}

		var existingUser userModels.User
		if err := db.Where("username = ?", login.Username).First(&existingUser).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{
				"error": "User not fount",
			})
			log.Printf("%s: user not found", op)
			return
		}

		if err := bcrypt.CompareHashAndPassword([]byte(existingUser.PasswordHash), []byte(login.Password)); err != nil {
			c.JSON(http.StatusUnauthorized, gin.H{
				"error": "Invalid password",
			})
			log.Printf("%s: invalid password", op)
			return
		}

		var existingToken tokenModels.UserToken

		if err := db.Where("user_id = ?", existingUser.ID).First(&existingToken).Error; err == nil {

			parsedToken, err := utils.ParseJWT(existingToken.Token, conf.Secret)
			if err != nil {
				c.JSON(http.StatusUnauthorized, gin.H{
					"error": "Failed to parse token",
				})
				log.Printf("%s: failed to parse token: %v", op, err)
				return
			}

			if !parsedToken.IsExpired() {
				c.JSON(http.StatusOK, gin.H{
					"message": "Your token is still valid",
					"token":   existingToken.Token,
				})
				log.Printf("%s: user %s already has a valid token", op, login.Username)
				return
			}
		} else if err != gorm.ErrRecordNotFound {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to check for existing token",
			})
			log.Printf("%s: failed to check for existing token: %v", op, err)
			return
		}

		tokenString, err := utils.GenerateJWT(existingUser.ID, conf.Secret)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to generate token",
			})
			log.Printf("%s: failed to generate token: %v", op, err)
			return
		}

		token := tokenModels.UserToken{
			UserID: existingUser.ID,
			Token:  tokenString,
		}

		if err := db.Create(&token).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"error": "Failed to save token",
			})
			log.Printf("%s: failed to save token: %v", op, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"message": "Your token has bin generated",
			"token":   tokenString,
		})

		log.Printf("%s: user %s logged in successfully", op, login.Username)
	}
}

func Me(db *gorm.DB) gin.HandlerFunc {
	op := "router.Me"

	return func(c *gin.Context) {
		user, status, err := utils.GetCurrentUser(db, c)
		if err != nil {
			c.JSON(status, gin.H{"error": err.Error()})
			return
		}

		userResponse := userSchemes.UserResponse{
			ID:        user.ID,
			Username:  user.Username,
			CreatedAt: user.CreatedAt.Format("15:04 02.01.2006"),
		}

		c.JSON(http.StatusOK, userResponse)
		log.Printf("%s: user %s retrieved successfully", op, user.Username)
	}
}

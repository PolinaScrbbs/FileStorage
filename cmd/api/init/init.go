package init

import (
	"FileStorage/api/middleware"
	"FileStorage/api/user/router"
	conf "FileStorage/config"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
)

func Run(r *gin.Engine, db *gorm.DB, config *conf.Config) {
	r.POST("/register", router.Register(db))
	r.POST("/login", router.Login(db, config))
	r.GET("/me", middleware.AuthMiddleware(config.Secret), router.Me(db))

	log.Printf("routers initialized")
}

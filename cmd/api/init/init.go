package init

import (
	"FileStorage/api/user/router"
	conf "FileStorage/config"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
)

func Run(r *gin.Engine, db *gorm.DB, confin *conf.Config) {
	r.POST("/register", router.Register(db))
	r.POST("login", router.Login(db, confin))

	log.Printf("routers initialized")
}

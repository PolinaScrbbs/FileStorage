package init

import (
	"FileStorage/api/user/router"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
)

func Run(r *gin.Engine, db *gorm.DB) {
	r.POST("/register", router.Register(db))

	log.Printf("routers initialized")
}

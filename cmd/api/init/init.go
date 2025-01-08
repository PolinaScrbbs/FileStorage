package init

import (
	fileRouter "FileStorage/api/file/router"
	"FileStorage/api/middleware"
	userRouter "FileStorage/api/user/router"
	conf "FileStorage/config"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"log"
)

func Run(r *gin.Engine, db *gorm.DB, conf *conf.Config) {
	r.POST("/register", userRouter.Register(db))
	r.POST("/login", userRouter.Login(db, conf))
	r.GET("/me", middleware.AuthMiddleware(conf.Secret), userRouter.Me(db))
	r.POST("/file/upload", middleware.AuthMiddleware(conf.Secret), fileRouter.AddFile(db, conf.Base_Save_Path))

	log.Printf("routers initialized")
}

package router

import (
	"FileStorage/api/file/schemes"
	"FileStorage/api/user/utils"
	fileModels "FileStorage/database/models/file"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"regexp"
)

func AddFile(db *gorm.DB, base_save_path string) gin.HandlerFunc {
	op := "router.AddFile"

	return func(c *gin.Context) {
		current_user, status, err := utils.GetCurrentUser(db, c)
		if err != nil {
			c.JSON(status, gin.H{"error": err.Error()})
			return
		}

		log.Printf("%s: User %s is uploading a file", op, current_user.Username)

		save_path := base_save_path + "/files" + "/" + current_user.Username

		if _, err := os.Stat(save_path); os.IsNotExist(err) {
			err = os.MkdirAll(save_path, os.ModePerm)
			if err != nil {
				log.Fatalf("%s: failed to create base directory: %v", op, err)
			}
		}

		var req schemes.CreateFileRequest

		if err := c.ShouldBind(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
				"input": req,
			})
			log.Printf("%s: invalid input data from user %s: %v", op, current_user.Username, err)
			return
		}

		validName := regexp.MustCompile("^[a-zA-Z0-9_-]+$")
		if !validName.MatchString(req.Name) {
			c.JSON(http.StatusBadRequest, gin.H{
				"error": "File name contains invalid characters. Only alphanumeric characters, hyphens, and underscores are allowed.",
			})
			log.Printf("%s: user %s provided an invalid file name: %v", op, current_user.Username, req.Name)
			return
		}

		file, err := req.File.Open()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to open file"})
			log.Printf("%s: user %s error opening file: %v", op, current_user.Username, err)
			return
		}
		defer file.Close()

		fileExtension := filepath.Ext(req.File.Filename)
		filePath := save_path + "/" + req.Name + fileExtension

		out, err := os.Create(filePath)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save file"})
			log.Printf("%s: user %s error creating file: %v", op, current_user.Username, err)
			return
		}
		defer out.Close()

		_, err = io.Copy(out, file)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to save file"})
			log.Printf("%s: user %s error copying file: %v", op, current_user.Username, err)
			return
		}

		log.Printf("%s: user %s uploaded file %s successfully", op, current_user.Username, req.File.Filename)

		new_file := fileModels.File{
			Name:      req.Name,
			Extension: fileExtension,
			Path:      filePath,
			CreatorID: current_user.ID,
		}

		if err := db.Create(&new_file).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save file"})
			log.Printf("%s: user %s failed to save file: %v", op, current_user.Username, err)
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "File uploaded and saved successfully"})
		log.Printf("%s: user %s saved file %s successfully", op, current_user.Username, req.File.Filename)
	}
}

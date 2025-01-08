package schemes

import "mime/multipart"

type CreateFileRequest struct {
	Name string                `form:"name" binding:"required"`
	File *multipart.FileHeader `form:"file" binding:"required"`
}

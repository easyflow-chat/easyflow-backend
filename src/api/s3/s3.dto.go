package s3

import "mime/multipart"

type UploadFileRequest struct {
	Upload *multipart.File `form:"upload" validate:"required`
}

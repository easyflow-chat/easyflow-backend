package s3

type GetDownloadURLRequest struct {
	BucketName string `json:"bucketName" validate:"required"`
	ObjectKey  string `json:"objectKey" validate:"required"`
}

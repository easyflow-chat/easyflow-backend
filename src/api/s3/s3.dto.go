package s3

type BucketRequest struct {
	Name string `json:"name" validate:"required"`
}

type GetDownloadURLRequest struct {
	BucketName string `json:"bucketName" validate:"required"`
	ObjectKey  string `json:"objectKey" validate:"required"`
}

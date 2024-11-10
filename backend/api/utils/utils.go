package utils

import (
	"easyflow-backend/api/s3"
	"easyflow-backend/common"

	"github.com/easyflow-chat/easyflow-backend/lib/database"

	"gorm.io/gorm"
)

func GenerateNewProfilePictureUrl(logger *common.Logger, cfg *common.Config, db *gorm.DB, user *database.User) {
	pictureUrl, err := s3.GenerateDownloadURL(logger, cfg, cfg.ProfilePictureBucketName, user.Id, 7*24*60*60)
	if err == nil {
		user.ProfilePicture = pictureUrl

		if err := db.Save(user).Error; err != nil {
			logger.PrintfWarning("Could not save the new ProfilePicture url for user: %s. Error: %s", user.Id, err)
		}
	}
}

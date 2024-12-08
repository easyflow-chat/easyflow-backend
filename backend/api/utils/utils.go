package utils

import (
	"easyflow-backend/api"
	"easyflow-backend/api/s3"
	"easyflow-backend/common"
	"easyflow-backend/enum"
	"encoding/json"
	"io"
	"net/http"
	"net/url"

	"github.com/easyflow-chat/easyflow-backend/lib/database"
	"github.com/easyflow-chat/easyflow-backend/lib/logger"

	"gorm.io/gorm"
)

func GenerateNewProfilePictureUrl(logger *logger.Logger, cfg *common.Config, db *gorm.DB, user *database.User) {
	pictureUrl, err := s3.GenerateDownloadURL(logger, cfg, cfg.ProfilePictureBucketName, user.ID, 7*24*60*60)
	if err == nil {
		user.ProfilePicture = pictureUrl

		if err := db.Save(user).Error; err != nil {
			logger.PrintfWarning("Could not save the new ProfilePicture url for user: %s. Error: %s", user.ID, err)
		}
	}
}

func CheckCloudflareTurnstile(logger *logger.Logger, cfg *common.Config, ip string, token string) (bool, *api.ApiError) {
	formData := url.Values{}
	formData.Add("secret", cfg.TurnstileSecret)
	formData.Add("response", token)
	formData.Add("remoteip", ip)

	res, err := http.PostForm(cfg.TurnstileUrl, formData)
	if err != nil {
		logger.PrintfError("Error verifying turnstile token: %s", err)
		return false, &api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		}
	}

	defer res.Body.Close()
	body, err := io.ReadAll(res.Body)
	if err != nil {
		logger.PrintfError("Error reading turnstile response: %s", err)
		return false, &api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		}
	}
	var jsonBody CloudflareTurnstileResponse
	err = json.Unmarshal(body, &jsonBody)
	if err != nil {
		logger.PrintfError("Error unmarshalling turnstile response: %s", err)
		return false, &api.ApiError{
			Code:  http.StatusInternalServerError,
			Error: enum.ApiError,
		}
	}

	logger.PrintfDebug("Action: %s", jsonBody.Action)

	if !jsonBody.Success {
		logger.PrintfWarning("Turnstile token verification failed: %s", jsonBody.ErrorCodes)
		return false, &api.ApiError{
			Code:    http.StatusUnauthorized,
			Error:   enum.InvalidTurnstile,
			Details: "Failed to validate the provided Turnstile token",
		}
	}

	return true, nil
}

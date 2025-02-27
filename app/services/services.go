package services

import (
	// "app/api/schema"
	"app/models"
	"app/repositories"

	"fmt"
	"errors"
)

func FindUser(userID int) (models.User, error) {
	return repositories.FindUser(userID)
}

func FindUserFavorites(userID int) []models.UserFavorite {
	return repositories.FindUserFavorites(userID)
}

func FindChartAsset(assetID int) (models.Chart, error) {
	return repositories.FindChartAsset(assetID)
}

func FindInsightAsset(assetID int) (models.Insight, error) {
	return repositories.FindInsightAsset(assetID)
}

func FindAudienceAsset(assetID int) (models.Audience, error) {
	return repositories.FindAudienceAsset(assetID)
}

func FindAllAssets() []models.Asset {
	return repositories.FindAllAssets()
}

func FindAsset(assetID int) (models.Asset, error) {
	return repositories.FindAsset(assetID)
}

func UpdateAsset(asset models.Asset) models.Asset {
	asset, err := repositories.SaveAsset(asset)
	if err != nil {
		panic(fmt.Sprintf("failed to update asset %d: cannot find asset", asset.ID))
	}
	return asset
}

// userFavoriteService
func CreateFavoriteAsset(assetID int, userID int) error {
	user, err := FindUser(userID)
	if err != nil {
		fmt.Printf("user %d not found", userID)
		return errors.New("user not found")
	}
	userFavorites := repositories.FindUserFavorites(userID)

	for _, uf := range userFavorites {
		if uf.Asset.ID == assetID {
			errorMessage := fmt.Sprintf("asset %d is already a favorite for user %d", assetID, userID)
			return errors.New(errorMessage)
		}
	}
	asset, err := FindAsset(assetID)
	if err != nil {
		fmt.Printf("asset %d not found", assetID)
		return errors.New("asset not found")
	}
	repositories.CreateUserFavorite(user, asset)
	return nil
}

func DeleteUserFavorite(userFavoriteID int) error {
	return repositories.DeleteUserFavorite(userFavoriteID)
}

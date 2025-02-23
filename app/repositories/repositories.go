package repositories

import (
	"app/models"
	"app/models/constants"

	"errors"
	"fmt"
	"slices"
)

var a1 = models.Asset{ID: 1, Description: "chart-1", Type: constants.ASSET_TYPE_CHART}
var a2 = models.Asset{ID: 2, Description: "chart-2", Type: constants.ASSET_TYPE_CHART}
var a3 = models.Asset{ID: 3, Description: "insight-1", Type: constants.ASSET_TYPE_INSIGHT}
var a4 = models.Asset{ID: 4, Description: "audience-1", Type: constants.ASSET_TYPE_AUDIENCE}
var a5 = models.Asset{ID: 5, Description: "audience-2", Type: constants.ASSET_TYPE_AUDIENCE}
// table assets
var assets = []models.Asset{a1, a2, a3, a4, a5}

// table charts
var charts = []models.Chart{
	{Asset: a1, Title: "chart-1-title", AxesTitles: "chart-1-AxesTitles", PlotData: "chart-1-PlotData"},
	{Asset: a2, Title: "chart-2-title", AxesTitles: "chart-2-AxesTitles", PlotData: "chart-2-PlotData"},
}

// table insights
var insights = []models.Insight{{Asset: a3, Text: "some nice insight"}}

var aud1 = models.Audience{Asset: a4, ID: 1}
var aud2 = models.Audience{Asset: a5, ID: 2}
// table audiences
var audiences = []models.Audience{aud1, aud2}

var ac1 = models.AudienceCharacteristic{ID: 1, AudienceID: 1, CharacteristicID: 1, CharacteristicValue: 1}
var ac2 = models.AudienceCharacteristic{ID: 2, AudienceID: 2, CharacteristicID: 1, CharacteristicValue: 2}
var ac3 = models.AudienceCharacteristic{ID: 3, AudienceID: 2, CharacteristicID: 2, CharacteristicValue: 300}
var ac4 = models.AudienceCharacteristic{ID: 4, AudienceID: 2, CharacteristicID: 4, CharacteristicValue: 4}
var ac5 = models.AudienceCharacteristic{ID: 5, AudienceID: 2, CharacteristicID: 5, CharacteristicValue: 4200}
var audienceCharacteristics = []models.AudienceCharacteristic{ac1, ac2, ac3, ac4, ac5}

var john = models.User{ID: 1, Name: "John Doe"}
var leroy = models.User{ID: 2, Name: "Leroy Jenkins"}
var users = []models.User {john, leroy}

// table userFavorites
var userFavorites = []models.UserFavorite{
	// passing references because in the database only the IDs are supposed to be saved
	{ID: 1, User: &john, Asset: &a1},
	{ID: 2, User: &john, Asset: &a3},
	{ID: 3, User: &leroy, Asset: &a2},
	{ID: 4, User: &leroy, Asset: &a5},
}

/* service methods */

func FindUser(userID int) (user models.User, Err error) {
	for _, u := range users {
		if u.ID == userID {
			return u, nil
		}
	}
	return models.User{}, errors.New("user not found")
}

// userFavoriteRepository
func FindUserFavorites(userID int) []models.UserFavorite {
	userFavoriteAssets := make([]models.UserFavorite, 0)
	for _, uf := range userFavorites {
		if uf.User.ID == userID {
			userFavoriteAssets = append(userFavoriteAssets, uf)
		}
	}
	return userFavoriteAssets
}

func CreateUserFavorite(user models.User, asset models.Asset) {
	nextID := len(userFavorites) + 1
	userFavorites = append(userFavorites, models.UserFavorite{ID: nextID, User: &user, Asset: &asset})
}

func DeleteUserFavorite(userFavoriteID int) error {
	index := -1
	for i, uf := range userFavorites {
		if uf.ID == userFavoriteID {
			index = i
			break
		}
	}
	if index < 0 {
		fmt.Printf("user favorite %d not found for deletion\n", userFavoriteID)
		return errors.New("user favorite not found for deletion")
	}
	userFavorites = slices.Delete(userFavorites, index, index + 1)
	return nil
}

func FindChartAsset(assetID int) (models.Chart, error) {
	for _, c := range charts {
		if c.ID == assetID {
			return c, nil
		}
	}
	return models.Chart{}, errors.New("chart not found")
}

func FindInsightAsset(assetID int) (models.Insight, error) {
	for _, i := range insights {
		if i.ID == assetID {
			return i, nil
		}
	}
	return models.Insight{}, errors.New("insight not found")
}

func FindAudienceAsset(assetID int) (models.Audience, error) {
	for _, a := range audiences {
		if a.Asset.ID == assetID {
			return a, nil
		}
	}
	return models.Audience{}, errors.New("audience not found")
}

func FindAudienceCharacteristics(audienceID int) ([]models.AudienceCharacteristic) {
	acs := make([]models.AudienceCharacteristic, 0)
	for _, ac := range audienceCharacteristics {
		if ac.AudienceID == audienceID {
			acs = append(acs, ac)
		}
	}
	return acs
}

func FindAllAssets() []models.Asset {
	return assets;
}

func FindAsset(assetID int) (models.Asset, error) {
	for _, a := range assets {
		if a.ID == assetID {
			return a, nil
		}
	}
	return models.Asset{}, errors.New("asset not found")
}

func SaveAsset(asset models.Asset) (models.Asset, error) {
	index := findAssetIndex(asset.ID)
	if index == -1 {
		return models.Asset{}, errors.New("asset not found")
	}
	assets[index] = asset
	return assets[index], nil
}

// when adding a database this will be removed
func findAssetIndex(assetID int) int {
	for index, a := range assets {
		if a.ID == assetID {
			return index
		}
	}
	return -1
}
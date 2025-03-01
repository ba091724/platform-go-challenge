package services

import (
	"app/api/schema"
	"app/services"
	
	"fmt"
	"net/http"
	"log"
)


func GetAssets(filter schema.AssetFilter) []schema.AssetDetailsDto {
	return services.FindAssets(filter)
}

func getAsset(id string) (schema.AssetDetailsDto, error) {
	dto, err := services.FindAsset(id)
	if err != nil {
		return schema.AssetDetailsDto{}, schema.NewApiError(http.StatusNotFound, err)
	}
	return dto, nil
}

func UpdateAsset(assetID string, updateRequest schema.AssetUpdateRequest) error {
	return services.UpdateAsset(assetID, updateRequest.Description)
}

func GetUserFavorites(userID string) []schema.UserFavoriteDto {
	ufs := services.FindUserFavorites(userID)
	if len(ufs) == 0 {
		return make([]schema.UserFavoriteDto, 0)
	}
	userFavoriteDtos := make([]schema.UserFavoriteDto, 0, len(ufs))
	for _, uf := range ufs {
		assetDto, err := getAsset(uf.AssetID)
		if err != nil {
			log.Printf("[!] failed to get asset details for asset %s, skipping...\n", uf.AssetID)
			continue
		} else {
		dto := schema.UserFavoriteDto{ID: uf.ID.Hex(), Details: assetDto}
		userFavoriteDtos = append(userFavoriteDtos, dto)
		}
	}
	return userFavoriteDtos
}

func CreateUserFavorite(userID string, assetID string) (string, error) {
	_, err := services.FindUser(userID)
	if err != nil {
		fmt.Printf("[X] user %s not found\n", userID)
		return "", schema.NewApiError(http.StatusNotFound, err)
	}
	if _, err := services.FindAsset(assetID); err != nil {
		fmt.Printf("[X] asset %s does not exist\n", assetID)
		return "", schema.NewApiError(http.StatusNotFound, err)
	}
	id, err := services.CreateFavoriteAsset(assetID, userID)
	if err != nil {
		return "", err
	}
	return id, nil
}

func DeleteUserFavorite(userID string, userFavoriteID string) error {
	_, err := services.FindUser(userID)
	if err != nil {
		fmt.Printf("user %s not found", userID)
		return schema.NewApiError(http.StatusNotFound, err)
	}
	return services.DeleteUserFavorite(userFavoriteID)
}

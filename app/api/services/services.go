package services

import (
	"net/http"

	"github.com/biter777/countries"

	"app/api/schema"
	"app/models"
	"app/models/constants"
	"app/repositories"
	"app/services"

	"errors"
	"fmt"
	// "http"
	"strconv"
	"strings"
)

func GetAssets() []schema.AssetDetailsDto {
	assets := services.FindAllAssets()
	if len(assets) == 0 {
		return make([]schema.AssetDetailsDto, 0)
	}
	response := make([]schema.AssetDetailsDto, 0, len(assets))
	for _, a := range assets {
		dto, err := getAssetDetails(a)
		if err != nil {
			fmt.Printf("[!] failed to get asset details for asset %d, skipping...\n", a.ID)
			continue
		} else {
			response = append(response, dto)
		}
	}
	return response
}

func UpdateAsset(assetID int, updateRequest schema.AssetUpdateRequest) (schema.AssetDetailsDto, error) {
	asset, err := services.FindAsset(assetID)
	if err != nil {
		//TODO return and handle bad request error (400)
		return schema.AssetDetailsDto{}, err
	}
	asset.Description = updateRequest.Description
	asset = services.UpdateAsset(asset)
	return getAssetDetails(asset)
}

func GetUserFavorites(userID int) []schema.UserFavoriteDto {
	ufs := services.FindUserFavorites(userID)
	if len(ufs) == 0 {
		return make([]schema.UserFavoriteDto, 0)
	}
	response := make([]schema.UserFavoriteDto, 0, len(ufs))
	for _, uf := range ufs {
		assetDto, err := getAssetDetails(*uf.Asset)
		if err != nil {
			fmt.Printf("[!] failed to get asset details for asset %s, skipping...\n", uf.Asset.ID)
			continue
		} else {
			responseDto := schema.UserFavoriteDto{ID: uf.ID, Details: assetDto}
			response = append(response, responseDto)
		}
	}
	return response
}

func getAssetDetails(asset models.Asset) (schema.AssetDetailsDto, error) {
	switch asset.Type {
	case constants.ASSET_TYPE_CHART:
		dto, err := getChartAsset(asset.ID)
		if err != nil {
			fmt.Printf("[X] chart asset %d not found\n", asset.ID)
			return schema.AssetDetailsDto{}, schema.NewApiError(http.StatusNotFound, err)
		}
		return schema.AssetDetailsDto{
			Asset:        schema.AssetDto{ID: asset.ID, Description: asset.Description},
			ChartDetails: &dto,
		}, nil
	case constants.ASSET_TYPE_INSIGHT:
		dto, err := getInsightAsset(asset.ID)
		if err != nil {
			fmt.Printf("[X] insight asset %d not found\n", asset.ID)
			return schema.AssetDetailsDto{}, schema.NewApiError(http.StatusNotFound, err)
		}
		return schema.AssetDetailsDto{
			Asset:          schema.AssetDto{ID: asset.ID, Description: asset.Description},
			InsightDetails: &dto,
		}, nil
	case constants.ASSET_TYPE_AUDIENCE:
		dto, err := getAudienceAsset(asset.ID)
		if err != nil {
			fmt.Printf("[X] audience asset %d not found\n", asset.ID)
			return schema.AssetDetailsDto{}, schema.NewApiError(http.StatusNotFound, err)
		}
		return schema.AssetDetailsDto{
			Asset:           schema.AssetDto{ID: asset.ID, Description: asset.Description},
			AudienceDetails: &dto,
		}, nil
	default:
		//TODO error log
		fmt.Println("[X] mapping not implemented yet for asset ", asset.ID)
		return schema.AssetDetailsDto{}, schema.NewApiError(http.StatusInternalServerError, nil)
	}
}

func CreateUserFavorite(userID int, assetID int) error {
	_, err := services.FindUser(userID)
	if err != nil {
		fmt.Printf("[X] user %d not found", userID)
		return schema.NewApiError(http.StatusNotFound, err)
	}
	if _, err := services.FindAsset(assetID); err != nil {
		fmt.Printf("[X] asset %d does not exist\n", assetID)
		return schema.NewApiError(http.StatusNotFound, err)
	}
	if err := services.CreateFavoriteAsset(assetID, userID); err != nil {
		// not good because might be a conflict as well
		// maybe customize those errors?
		return schema.NewApiError(http.StatusBadRequest, err)
	}
	return nil
}

func DeleteUserFavorite(userID int, userFavoriteID int) error {
	_, err := services.FindUser(userID)
	if err != nil {
		fmt.Printf("user %d not found", userID)
		return schema.NewApiError(http.StatusNotFound, err)
	}
	return services.DeleteUserFavorite(userFavoriteID)
}

func getChartAsset(assetID int) (schema.ChartAssetDto, error) {
	a, err := services.FindChartAsset(assetID)
	if err != nil {
		return schema.ChartAssetDto{}, schema.NewApiError(http.StatusNotFound, err)
	}
	return schema.ChartAssetDto{
		Title:      a.Title,
		AxesTitles: a.AxesTitles,
		PlotData:   a.PlotData,
	}, nil
}

func getInsightAsset(assetID int) (schema.InsightAssetDto, error) {
	a, err := services.FindInsightAsset(assetID)
	if err != nil {
		return schema.InsightAssetDto{}, schema.NewApiError(http.StatusNotFound, err)
	}
	return schema.InsightAssetDto{Text: a.Text}, nil
}

// this could be placed in a mapping method, or mapped on an entity property (?)
func getAudienceAsset(assetID int) (schema.AudienceAssetDto, error) {
	audience, err := services.FindAudienceAsset(assetID)
	if err != nil {
		return schema.AudienceAssetDto{}, schema.NewApiError(http.StatusNotFound, err)
	}
	var characteristics []string
	for _, ac := range repositories.FindAudienceCharacteristics(audience.ID) {
		str, err := getCharacteristicStr(ac.CharacteristicID, ac.CharacteristicValue)
		if err != nil {
			fmt.Printf("[!] something went wrong when resolving characteristic with id=%d, skipping...\n", ac.CharacteristicID)
			continue
		} else {
			characteristics = append(characteristics, str)
		}
	}
	return schema.AudienceAssetDto{Characteristics: characteristics}, nil
}

func getCharacteristicStr(characteristicID int, characteristicValue int) (string, error) {
	var sb strings.Builder
	if characteristicID == constants.CHARACTERISTIC_AGE_GROUP {
		ageGroup := getAgeGroupText(characteristicValue)
		sb.WriteString(schema.CHARACTERISTIC_AGE_GROUP)
		sb.WriteString(": ")
		sb.WriteString(ageGroup)
	} else if characteristicID == constants.CHARACTERISTIC_BIRTH_COUNTRY {
		country := countries.ByNumeric(characteristicValue).Info().Name
		sb.WriteString(schema.CHARACTERISTIC_BIRTH_COUNTRY)
		sb.WriteString(": ")
		sb.WriteString(country)
	} else if characteristicID == constants.CHARACTERISTIC_GENDER {
		gender := getGenderText(characteristicValue)
		sb.WriteString(schema.CHARACTERISTIC_GENDER)
		sb.WriteString(": ")
		sb.WriteString(gender)
	} else if characteristicID == constants.CHARACTERISTIC_PURCHASES_LAST_MONTH {
		sb.WriteString(schema.CHARACTERISTIC_PURCHASES_LAST_MONTH)
		sb.WriteString(": ")
		sb.WriteString(strconv.Itoa(characteristicValue))
	} else if characteristicID == constants.CHARACTERISTIC_SOCIAL_MEDIA_DAILY_HOURS {
		sb.WriteString(schema.CHARACTERISTIC_SOCIAL_MEDIA_DAILY_HOURS)
		sb.WriteString(": ")
		sb.WriteString(strconv.Itoa(characteristicValue))
	} else {
		return "", errors.New("unhandled characteristic id")
	}
	return sb.String(), nil
}

func getGenderText(genderID int) string {
	if genderID == 1 {
		return schema.GENDER_FEMALE
	}
	return schema.GENDER_MALE
}

func getAgeGroupText(ageGroupID int) string {
	switch ageGroupID {
	case constants.AGE_GROUP_17_19:
		return "17 to 19"
	case constants.AGE_GROUP_20_29:
		return "20 to 29"
	case constants.AGE_GROUP_30_39:
		return "30 to 39"
	case constants.AGE_GROUP_40_49:
		return "40 to 49"
	case constants.AGE_GROUP_50_59:
		return "50 to 59"
	case constants.AGE_GROUP_60_69:
		return "60 to 69"
	case constants.AGE_GROUP_70_PLUS:
		return "70+"
	default:
		return "UNKNOWN"
	}
}

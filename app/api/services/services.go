package services

import (
	"errors"

	"github.com/biter777/countries"

	"app/api/payload"
	"app/models/constants"
	"app/repositories"
	"app/services"
	"fmt"
	"strconv"
	"strings"
)

func GetUserFavorites(userID int) []payload.UserFavoriteDto {
	ufs := services.FindUserFavorites(userID)
	if len(ufs) == 0 {
		return make([]payload.UserFavoriteDto, 0)
	}
	response := make([]payload.UserFavoriteDto, 0, len(ufs))
	for _, uf := range ufs {
	SWITCH:
		switch {
		case uf.Asset.Type == constants.ASSET_TYPE_CHART:
			dto, err := getChartAsset(uf.Asset.ID)
			if err != nil {
				fmt.Printf("[X] chart asset %d not found\n", uf.Asset.ID)
				continue
			}
			assetDetailsDto := payload.UserFavoriteDto{
				ID: uf.ID,
				Details: payload.AssetDetailsDto{
					Asset:        payload.AssetDto{ID: uf.Asset.ID},
					ChartDetails: &dto,
				},
			}
			response = append(response, assetDetailsDto)
			break SWITCH
		case uf.Asset.Type == constants.ASSET_TYPE_INSIGHT:
			dto, err := getInsightAsset(uf.Asset.ID)
			if err != nil {
				fmt.Printf("[X] insight asset %d not found\n", uf.Asset.ID)
				continue
			}
			assetDetailsDto := payload.UserFavoriteDto{
				ID: uf.ID,
				Details: payload.AssetDetailsDto{
					Asset:          payload.AssetDto{ID: uf.Asset.ID},
					InsightDetails: &dto,
				},
			}
			response = append(response, assetDetailsDto)
			break SWITCH
		case uf.Asset.Type == constants.ASSET_TYPE_AUDIENCE:
			dto, err := getAudienceAsset(uf.Asset.ID)
			if err != nil {
				fmt.Printf("[X] audience asset %d not found\n", uf.Asset.ID)
				continue
			}
			assetDetailsDto := payload.UserFavoriteDto{
				ID: uf.ID,
				Details: payload.AssetDetailsDto{
					Asset:           payload.AssetDto{ID: uf.Asset.ID},
					AudienceDetails: &dto,
				},
			}
			response = append(response, assetDetailsDto)
			break SWITCH
		default:
			fmt.Println("mapping not implemented yet for asset ", uf.Asset.ID)
		}
	}
	return response
}

func CreateUserFavorite(userID int, assetID int) error {
	_, err := services.FindUser(userID)
	if err != nil {
		fmt.Printf("user %d not found", userID)
		return err
	}
	if !services.AssetExists(assetID) {
		fmt.Printf("[X] asset %d does not exist\n", assetID)
		return errors.New("unknown asset")
	}
	services.CreateFavoriteAsset(assetID, userID)
	return nil
}

func DeleteUserFavorite(userID int, userFavoriteID int) error {
	_, err := services.FindUser(userID)
	if err != nil {
		fmt.Printf("user %d not found", userID)
		return err
	}
	return services.DeleteUserFavorite(userFavoriteID)
}

func getChartAsset(assetID int) (payload.ChartAssetDto, error) {
	a, err := services.FindChartAsset(assetID)
	if err != nil {
		return payload.ChartAssetDto{}, err
	}
	return payload.ChartAssetDto{
		Title:      a.Title,
		AxesTitles: a.AxesTitles,
		PlotData:   a.PlotData,
	}, nil
}

func getInsightAsset(assetID int) (payload.InsightAssetDto, error) {
	a, err := services.FindInsightAsset(assetID)
	if err != nil {
		return payload.InsightAssetDto{}, err
	}
	return payload.InsightAssetDto{Text: a.Text}, nil
}

// this could be placed in a mapping method, or mapped on an entity property (?)
func getAudienceAsset(assetID int) (payload.AudienceAssetDto, error) {
	audience, err := services.FindAudienceAsset(assetID)
	if err != nil {
		return payload.AudienceAssetDto{}, err
	}
	var characteristics []string
	for _, ac := range repositories.FindAudienceCharacteristics(audience.ID) {
		str := getCharacteristicStr(ac.CharacteristicID, ac.CharacteristicValue)
		characteristics = append(characteristics, str)
	}
	return payload.AudienceAssetDto{Characteristics: characteristics}, nil
}

func getCharacteristicStr(characteristicID int, characteristicValue int) string {
	var sb strings.Builder
	if characteristicID == constants.CHARACTERISTIC_AGE_GROUP {
		ageGroup := getAgeGroupText(characteristicValue)
		sb.WriteString(payload.CHARACTERISTIC_AGE_GROUP)
		sb.WriteString(": ")
		sb.WriteString(ageGroup)
	} else if characteristicID == constants.CHARACTERISTIC_BIRTH_COUNTRY {
		country := countries.ByNumeric(characteristicValue).Info().Name
		sb.WriteString(payload.CHARACTERISTIC_BIRTH_COUNTRY)
		sb.WriteString(": ")
		sb.WriteString(country)
	} else if characteristicID == constants.CHARACTERISTIC_GENDER {
		gender := getGenderText(characteristicValue)
		sb.WriteString(payload.CHARACTERISTIC_GENDER)
		sb.WriteString(": ")
		sb.WriteString(gender)
	} else if characteristicID == constants.CHARACTERISTIC_PURCHASES_LAST_MONTH {
		sb.WriteString(payload.CHARACTERISTIC_PURCHASES_LAST_MONTH)
		sb.WriteString(": ")
		sb.WriteString(strconv.Itoa(characteristicValue))
	} else if characteristicID == constants.CHARACTERISTIC_SOCIAL_MEDIA_DAILY_HOURS {
		sb.WriteString(payload.CHARACTERISTIC_SOCIAL_MEDIA_DAILY_HOURS)
		sb.WriteString(": ")
		sb.WriteString(strconv.Itoa(characteristicValue))
	} else {
		//TODO default handling needed
	}
	return sb.String()
}

func getGenderText(genderID int) string {
	if genderID == 1 {
		return payload.GENDER_FEMALE
	}
	return payload.GENDER_MALE
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

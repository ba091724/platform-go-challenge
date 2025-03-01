package repositories

import (
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"app/api/schema"
	"app/configs"
	"app/models"
	"app/models/constants"

	"fmt"
	"net/http"
	"errors"
	"context"
	"log"
)

/* collections */

var assetsCollection *mongo.Collection = configs.GetCollection(configs.DB, "assets")
var audienceCharacteristicsCollection *mongo.Collection = configs.GetCollection(configs.DB, "audienceCharacteristics")
var userFavoritesCollection *mongo.Collection = configs.GetCollection(configs.DB, "userFavorites")
var userCollection *mongo.Collection = configs.GetCollection(configs.DB, "users")

/* service methods */

func FindUser(userID string) (user models.User, Err error) {
	userId, errx := primitive.ObjectIDFromHex(userID)
	if errx != nil {
		return models.User{}, schema.NewApiError(http.StatusNotFound, errors.New("user not found"))
	}
	var result models.User
	err := userCollection.FindOne(
		context.TODO(),
		bson.D{{"_id", userId}},
	).Decode(&result)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in
		// the collection.
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.User{}, err
		}
		log.Panic(err)
	}
	return result, nil
}

func FindUserFavorites(userID string) []models.UserFavorite {
	cursor, err := userFavoritesCollection.Find(context.TODO(), bson.D{{"userId", userID}})
	if err != nil {
		log.Panic(err)
	}
	var results []models.UserFavorite
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Panic(err)
	}
	return results
}

func CreateUserFavorite(userID string, assetID string) string {
	res, err := userFavoritesCollection.InsertOne(context.TODO(), models.UserFavorite{UserID: userID, AssetID: assetID})
	if err != nil {
		log.Panicf("failed to insert new user favorite {userId=%s, assetId=%s}\n", userID, assetID)
	}
	return fmt.Sprintf("%s", res.InsertedID)
}

func DeleteUserFavorite(userFavoriteID string) error {
	userFavoriteId, errx := primitive.ObjectIDFromHex(userFavoriteID)
	if errx != nil {
		return schema.NewApiError(http.StatusNotFound, errors.New("user favorite not found"))
	}
	res, err := userFavoritesCollection.DeleteOne(context.TODO(), bson.M{"_id": userFavoriteId})
	if err != nil {
		return schema.NewApiError(http.StatusInternalServerError, err)
	}
	if res.DeletedCount == 0 {
		return schema.NewApiError(http.StatusNotFound, errors.New("user favorite not found"))
	}
	return nil
}

func FindAsset(assetID string) (models.AssetVO, error) {
	assetId, errx := primitive.ObjectIDFromHex(assetID)
	if errx != nil {
		return models.AssetVO{}, errx
	}
	var result bson.M
	err := assetsCollection.FindOne(
		context.TODO(),
		bson.D{{"_id", assetId}},
	).Decode(&result)
	if err != nil {
		// ErrNoDocuments means that the filter did not match any documents in
		// the collection.
		if errors.Is(err, mongo.ErrNoDocuments) {
			return models.AssetVO{}, err
		}
		log.Panic(err)
	}
	return getAssetVo(result), nil //TODO not good, must return error
}

func FindAssets(filter schema.AssetFilter) []models.AssetVO {
	bsonFilter := getBsonFilter(filter)
	cursor, err := assetsCollection.Find(context.TODO(), bsonFilter)
	if err != nil {
		log.Panic(err)
	}

	var results []bson.M
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Panic(err)
	}
	var assetVos []models.AssetVO = make([]models.AssetVO, len(results))
	for _, r := range results {
		assetVos = append(assetVos, getAssetVo(r))
	}
	return assetVos
}

func UpdateAsset(request models.AssetUpdateRequest) error {
	assetId, errx := primitive.ObjectIDFromHex(request.AssetID)
	if errx != nil {
		return schema.NewApiError(http.StatusNotFound, errors.New("asset not found"))
	}
	opts := options.FindOneAndUpdate().SetUpsert(true)
	filter := bson.D{{"_id", assetId}}
	update := bson.D{{"$set", bson.D{{"description", request.Description}}}}
	var updatedAsset models.Asset
	err := assetsCollection.FindOneAndUpdate(
		context.TODO(),
		filter,
		update,
		opts,
	).Decode(&updatedAsset)
	if err != nil {
		if errors.Is(err, mongo.ErrNoDocuments) {
			return schema.NewApiError(http.StatusNotFound, errors.New("asset not found"))
		}
		log.Print(err)
		return schema.NewApiError(http.StatusInternalServerError, err)
	}
	return nil
}

func getBsonFilter(filter schema.AssetFilter) bson.D {
	// not a good way because I have to cover all the property combinations of filter arg
	if filter.AssetID != "" {
		return bson.D{{"assetId", filter.AssetID}}
	}
	return bson.D{{}}
}

func getAssetVo(raw bson.M) models.AssetVO {
	assetType := int(raw["type"].(int32))
	assetId := raw["_id"].(primitive.ObjectID).Hex()
	description := raw["description"].(string)
	if assetType == constants.ASSET_TYPE_CHART {
		title := raw["description"].(string)
		axesTitles := raw["axesTitles"].(string)
		plotData := raw["plotData"].(string)
		return models.NewChartAsset(assetId, description, title, axesTitles, plotData)
	}
	if assetType == constants.ASSET_TYPE_INSIGHT {
		text := raw["text"].(string)
		return models.NewInsightAsset(assetId, description, text)
	}
	if assetType == constants.ASSET_TYPE_AUDIENCE {
		audienceCharacteristics := findAudienceCharacteristics(assetId)
		return models.NewAudienceAsset(assetId, description, audienceCharacteristics)
	}
	panic(fmt.Sprintf("failed to map raw bson '%s' to AssetVO, unknown asset type\n", raw))
}

func findAudienceCharacteristics(assetID string) []models.AudienceCharacteristic {
	cursor, err := audienceCharacteristicsCollection.Find(context.TODO(), bson.D{{"assetId", assetID}})
	if err != nil {
		log.Panic(err)
	}
	var results []models.AudienceCharacteristic
	if err = cursor.All(context.TODO(), &results); err != nil {
		log.Panic(err)
	}
	return results
}

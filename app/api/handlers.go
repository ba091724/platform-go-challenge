package api

import (
	"github.com/gin-gonic/gin"

	"app/api/schema"
	"app/api/services"

	"fmt"
	"net/http"
	// "reflect"
	"strconv"
	// "validator"
	// "json"
	// "ioutil"
	// "errors"
)

func GetAssets(c *gin.Context) {
	assets := services.GetAssets()
	c.JSON(http.StatusOK, assets)
}

// TODO also check if user exists, because a user will be firing this request. or make compelling arguments about that being taken care of the security framework
func UpdateAsset(c *gin.Context) {
	assetId := getIntParamFromPath("assetId", c)
	var json schema.AssetUpdateRequest
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	//TODO custom error handling
	asset, err := services.UpdateAsset(assetId, json)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	} else {
		c.JSON(http.StatusOK, asset)
	}
}

func GetUserFavorites(c *gin.Context) {
	userId := getIntParamFromPath("userId", c)
	userFavorites := services.GetUserFavorites(userId)
	c.JSON(http.StatusOK, userFavorites)
}

func CreateUserFavorite(c *gin.Context) {
	userId := getIntParamFromPath("userId", c)
	var json schema.UserFavoriteRequest
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := services.CreateUserFavorite(userId, json.AssetID)
	if err != nil {
		httpError, ok := err.(schema.HttpError)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(httpError.Status(), gin.H{"error": httpError.Error()})
	}
}

func DeleteUserFavorite(c *gin.Context) {
	userId := getIntParamFromPath("userId", c)
	userFavoriteId := getIntParamFromPath("favoriteId", c)
	err := services.DeleteUserFavorite(userId, userFavoriteId)
	if err != nil {
		httpError, ok := err.(schema.HttpError)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(httpError.Status(), gin.H{"error": httpError.Error()})
	}
}

func getIntParamFromPath(paramName string, c *gin.Context) int {
	defer func() {
        if r := recover(); r != nil {
            fmt.Println("failed to parse int param", r)
			c.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("Param '%s' is a number", paramName)})
			return
        }
    }()
	paramValue, err := strconv.Atoi(c.Param(paramName))
	if err != nil {
		// might be useful
		panic(schema.ValidationError{Key: paramName})
	}
	return paramValue
}
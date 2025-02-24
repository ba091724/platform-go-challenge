package api

import (
	"errors"

	"github.com/gin-gonic/gin"

	"app/api/payload"
	"app/api/services"

	"fmt"
	"net/http"
	"strconv"
)

func GetAssets(c *gin.Context) {
    assets := services.GetAssets()
    c.JSON(http.StatusOK, assets)
}

//TODO also check if user exists, because a user will be firing this request. or make compelling arguments about that being taken care of the security framework
func UpdateAsset(c *gin.Context) {
    assetId := getIntParamFromPath("assetId", c)
    var json payload.AssetUpdateRequest
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
    var json payload.UserFavoriteRequest
    if err := c.BindJSON(&json); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
        return
    }
    err := services.CreateUserFavorite(userId, json.AssetID)
    if err.Err != nil {
        httpErrorCode := getErrorCode(err)
        c.JSON(httpErrorCode, gin.H{"error": err.Error()})
    }
}

func DeleteUserFavorite(c *gin.Context) {
    userId := getIntParamFromPath("userId", c)
    userFavoriteId := getIntParamFromPath("favoriteId", c)
    services.DeleteUserFavorite(userId, userFavoriteId)
}

func getIntParamFromPath(paramName string, c *gin.Context) int {
    paramValue, err := strconv.Atoi(c.Param(paramName))
    if err != nil {
        fmt.Printf("[X] Path param '%s' is mandatory\n", paramName)
		panic(err)
	}
    return paramValue
}

func getErrorCode(err payload.ErrHttp) int {
    if errors.Is(err, payload.ErrBadRequest) {
        return http.StatusBadRequest
    }
    if errors.Is(err, payload.ErrNotFound) {
        return http.StatusNotFound
    }
    if errors.Is(err, payload.ErrConflict) {
        return http.StatusConflict
    }
    return http.StatusInternalServerError
}
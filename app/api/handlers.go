package api

import (
	"github.com/gin-gonic/gin"

	"app/api/services"
    "app/api/payload"

	"net/http"
	"strconv"
	"fmt"
)

func GetAssets(c *gin.Context) {
    assets := services.GetAssets()
    c.JSON(http.StatusOK, assets)
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
    fmt.Println(json)
    services.CreateUserFavorite(userId, json.AssetID)
}

func DeleteUserFavorite(c *gin.Context) {
    userId := getIntParamFromPath("userId", c)
    userFavoriteId := getIntParamFromPath("favoriteId", c)
    services.DeleteUserFavorite(userId, userFavoriteId)
}

func getIntParamFromPath(paramName string, c *gin.Context) int {
    // paramValue, err := strconv.Atoi(c.Request.PathValue(paramName))
    paramValue, err := strconv.Atoi(c.Param(paramName))
    if err != nil {
        fmt.Printf("[X] Path param '%s' is mandatory\n", paramName)
		panic(err)
	}
    return paramValue
}
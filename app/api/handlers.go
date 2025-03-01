package api

import (
	"github.com/gin-gonic/gin"

	"app/api/schema"
	"app/api/services"

	"net/http"
)

func GetAssets(c *gin.Context) {
	assets := services.GetAssets(schema.AssetFilter{})
	c.JSON(http.StatusOK, assets)
}

// TODO also check if user exists, because a user will be firing this request. or make compelling arguments about that being taken care of the security framework
func UpdateAsset(c *gin.Context) {
	assetId := c.Param("assetId")
	var json schema.AssetUpdateRequest
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	err := services.UpdateAsset(assetId, json)
	if err != nil {
		httpError, ok := err.(schema.HttpError)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(httpError.Status(), gin.H{"error": httpError.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{})
	}
}

func GetUserFavorites(c *gin.Context) {
	userId := c.Param("userId")
	userFavorites := services.GetUserFavorites(userId)
	c.JSON(http.StatusOK, userFavorites)
}

func CreateUserFavorite(c *gin.Context) {
	userId := c.Param("userId")
	var json schema.UserFavoriteRequest
	if err := c.BindJSON(&json); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	id, err := services.CreateUserFavorite(userId, json.AssetID)
	if err != nil {
		httpError, ok := err.(schema.HttpError)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		}
		c.JSON(httpError.Status(), gin.H{"error": httpError.Error()})
	} else {
		c.JSON(http.StatusOK, gin.H{"id": id}) // not good, must return schema.UserFavoriteDto
	}
}

func DeleteUserFavorite(c *gin.Context) {
	userId := c.Param("userId")
	userFavoriteId := c.Param("favoriteId")
	err := services.DeleteUserFavorite(userId, userFavoriteId)
	if err != nil {
		httpError, ok := err.(schema.HttpError)
		if !ok {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		} else {
			c.JSON(httpError.Status(), gin.H{"error": httpError.Error()})
		}
	}
}
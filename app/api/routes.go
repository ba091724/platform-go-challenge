package api

import "github.com/gin-gonic/gin"

func SetupRoutes(router *gin.Engine) {
    router.GET("/users/:userId/favorites", GetUserFavorites)
    router.POST("/users/:userId/favorites", CreateUserFavorite)
    router.DELETE("/users/:userId/favorites/:favoriteId", DeleteUserFavorite)
}
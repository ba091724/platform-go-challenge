package main

import (
    "app/api"
    "github.com/gin-gonic/gin"
)

func main() {
    router := gin.Default()
    api.SetupRoutes(router)
    router.Run(":8080")
}
package main

import (
    "github.com/gin-gonic/gin"
    
    "app/api"
    "app/configs"
)

func main() {
    configs.ConnectDB()
    router := gin.Default()
    api.SetupRoutes(router)
    router.Run(":8080")
}
package main

import (
    "github.com/gin-gonic/gin"
    
    "app/api"
)

func main() {
    router := gin.Default()
    api.SetupRoutes(router)
    router.Run(":8080")
}
package main

import (
  "github.com/gin-gonic/gin"
  "go-by-example/domain"
)

func main() {
  router := gin.Default()
  router.GET("/animal", domain.GetAnimals)
  router.GET("/animal/:id", domain.GetAnimalById)

  router.POST("/animal", domain.CreateAnimal)
  router.DELETE("/animal/:id", domain.DeleteAnimalById)
  router.PUT("/animal", domain.UpdateAnimal)

  router.Run("localhost:8080")
}



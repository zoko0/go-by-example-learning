package domain

import (
  "log"
  "github.com/gin-gonic/gin"
  "net/http"
  "go-by-example/model"
)

func GetAnimals(c *gin.Context) {
  animals, err := getAllAnimals()
  if err != nil {
    c.IndentedJSON(http.StatusBadRequest, gin.H{"errorMessage": err})
    return
  }

  c.IndentedJSON(http.StatusOK, animals)
}

func GetAnimalById(c *gin.Context) {
  id := c.Param("id")
  if id == "" {
    c.IndentedJSON(http.StatusBadRequest, gin.H{"errorMessage": "Missing ID path parameter"})
    return
  }

  animal, err := getAnimalById(id)
  if err != nil {
    c.IndentedJSON(http.StatusNotFound, gin.H{"errorMessage": "Animal with given ID not found"})
  }

  c.IndentedJSON(http.StatusOK, animal)
}

func CreateAnimal(c *gin.Context) {

  var newAnimal model.Animal

  if err := c.BindJSON(&newAnimal); err != nil {
    return
  }

  createdAnimal, err := addAnimal(newAnimal)
  if err != nil {
    c.IndentedJSON(http.StatusBadRequest, gin.H{"errorMessage": err})
    return
  }

  c.IndentedJSON(http.StatusCreated, createdAnimal)
}

func UpdateAnimal(c *gin.Context) {
  var newAnimal model.Animal

  if jsonErr := c.BindJSON(&newAnimal); jsonErr != nil {
    log.Println("error")
    c.IndentedJSON(http.StatusBadRequest, gin.H{"errorMessage": "Incorrect input data"})
    return
  }
  
  updatedAnimal, err := updateAnimal(newAnimal)
  if err != nil {
    c.IndentedJSON(http.StatusBadRequest, gin.H{"errorMessage": err})
    return
  }

  c.IndentedJSON(http.StatusOK, updatedAnimal)
}

func DeleteAnimalById(c *gin.Context) {
  id := c.Param("id")
  if id == "" {
     c.IndentedJSON(http.StatusBadRequest, gin.H{"errorMessage": "Missing ID in path"})
     return
  }
  _, err := deleteAnimal(id)
  if err != nil {
    c.IndentedJSON(http.StatusBadRequest, gin.H{"errorMessage": err})
    return
  }

  c.IndentedJSON(http.StatusOK, gin.H{"message": "Animal deleted"})
}

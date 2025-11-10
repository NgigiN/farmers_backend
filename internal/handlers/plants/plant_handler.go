package plants

import (
	plants "farm-backend/internal/models/plants"
	services "farm-backend/internal/services/plants"
	"strconv"

	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PlantHandler struct {
	PlantService *services.PlantService
}

func NewPlantHandler(plantService *services.PlantService) *PlantHandler {
	return &PlantHandler{PlantService: plantService}
}

func (h *PlantHandler) CreatePlant(c *gin.Context) {
	UserID := c.GetUint("user_id")
	var plant plants.Plant
	if err := c.ShouldBindJSON(&plant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.PlantService.Create(UserID, &plant); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, plant)
}

func (h *PlantHandler) ListPlants(c *gin.Context) {
	UserID := c.GetUint("user_id")
	plants, err := h.PlantService.List(UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, plants)
}

func (h *PlantHandler) GetPlant(c *gin.Context) {
	UserID := c.GetUint("user_id")
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plant ID"})
		return
	}
	plant, err := h.PlantService.Get(UserID, uint(idUint))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Plant not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, plant)
}

func (h *PlantHandler) UpdatePlant(c *gin.Context) {
	UserID := c.GetUint("user_id")
	id := c.Param("id")
	var plant plants.Plant
	if err := c.ShouldBindJSON(&plant); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plant ID"})
		return
	}
	if err := h.PlantService.Update(UserID, uint(idUint), &plant); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Plant not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	updated, _ := h.PlantService.Get(UserID, uint(idUint))
	c.JSON(http.StatusOK, updated)
}

func (h *PlantHandler) DeletePlant(c *gin.Context) {
	UserID := c.GetUint("user_id")
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plant ID"})
		return
	}
	if err := h.PlantService.Delete(UserID, uint(idUint)); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Plant not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Plant deleted successfully"})
}

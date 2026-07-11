package plants

import (
	plants "farm-backend/internal/services/plants"
	"farm-backend/internal/validation"
	"strconv"

	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type PlantHandler struct {
	PlantService *plants.PlantService
}

func NewPlantHandler(plantService *plants.PlantService) *PlantHandler {
	return &PlantHandler{PlantService: plantService}
}

func (h *PlantHandler) CreatePlant(c *gin.Context) {
	UserID := c.GetUint("user_id")
	var req validation.PlantRequest
	if err := validation.BindAndValidate(c, &req); err != nil {
		validation.RespondBindingError(c, err)
		return
	}
	entity := validation.PlantFromRequest(&req)
	if err := h.PlantService.Create(UserID, entity); err != nil {
		validation.RespondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, entity)
}

func (h *PlantHandler) ListPlants(c *gin.Context) {
	UserID := c.GetUint("user_id")
	plants, err := h.PlantService.List(UserID)
	if err != nil {
		validation.RespondError(c, err)
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
			validation.RespondNotFound(c, "Plant not found")
			return
		}
		validation.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, plant)
}

func (h *PlantHandler) UpdatePlant(c *gin.Context) {
	UserID := c.GetUint("user_id")
	id := c.Param("id")
	var req validation.PlantRequest
	if err := validation.BindAndValidate(c, &req); err != nil {
		validation.RespondBindingError(c, err)
		return
	}
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid plant ID"})
		return
	}
	entity := validation.PlantFromRequest(&req)
	if err := h.PlantService.Update(UserID, uint(idUint), entity); err != nil {
		if err == gorm.ErrRecordNotFound {
			validation.RespondNotFound(c, "Plant not found")
			return
		}
		validation.RespondError(c, err)
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
			validation.RespondNotFound(c, "Plant not found")
			return
		}
		validation.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Plant deleted successfully"})
}
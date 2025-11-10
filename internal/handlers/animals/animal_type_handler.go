package animals

import (
	animals "farm-backend/internal/models/animals"
	services "farm-backend/internal/services/animals"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AnimalTypeHandler struct {
	AnimalTypeService *services.AnimalTypeService
}

func NewAnimalTypeHandler(animalTypeService *services.AnimalTypeService) *AnimalTypeHandler {
	return &AnimalTypeHandler{AnimalTypeService: animalTypeService}
}

func (h *AnimalTypeHandler) CreateAnimalType(c *gin.Context) {
	UserID := c.GetUint("user_id")
	var animalType animals.AnimalType
	if err := c.ShouldBindJSON(&animalType); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.AnimalTypeService.Create(UserID, &animalType); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, animalType)
}

func (h *AnimalTypeHandler) ListAnimalTypes(c *gin.Context) {
	UserID := c.GetUint("user_id")
	animalTypes, err := h.AnimalTypeService.List(UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, animalTypes)
}

func (h *AnimalTypeHandler) GetAnimalType(c *gin.Context) {
	UserID := c.GetUint("user_id")
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid animal type ID"})
		return
	}
	animalType, err := h.AnimalTypeService.Get(UserID, uint(idUint))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Animal type not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, animalType)
}

func (h *AnimalTypeHandler) UpdateAnimalType(c *gin.Context) {
	UserID := c.GetUint("user_id")
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid animal type ID"})
		return
	}
	var animalType animals.AnimalType
	if err := c.ShouldBindJSON(&animalType); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.AnimalTypeService.Update(UserID, uint(idUint), &animalType); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Animal type not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	updated, _ := h.AnimalTypeService.Get(UserID, uint(idUint))
	c.JSON(http.StatusOK, updated)
}

func (h *AnimalTypeHandler) DeleteAnimalType(c *gin.Context) {
	UserID := c.GetUint("user_id")
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid animal type ID"})
		return
	}
	if err := h.AnimalTypeService.Delete(UserID, uint(idUint)); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Animal type not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Animal type deleted successfully"})
}

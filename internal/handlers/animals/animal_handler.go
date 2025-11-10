package animals

import (
	animals "farm-backend/internal/models/animals"
	services "farm-backend/internal/services/animals"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type AnimalHandler struct {
	AnimalService *services.AnimalService
}

func NewAnimalHandler(animalService *services.AnimalService) *AnimalHandler {
	return &AnimalHandler{AnimalService: animalService}
}

func (h *AnimalHandler) CreateAnimal(c *gin.Context) {
	UserID := c.GetUint("user_id")
	var animal animals.Animal
	if err := c.ShouldBindJSON(&animal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.AnimalService.Create(UserID, &animal); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, animal)
}

func (h *AnimalHandler) ListAnimals(c *gin.Context) {
	UserID := c.GetUint("user_id")
	animals, err := h.AnimalService.List(UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, animals)
}

func (h *AnimalHandler) GetAnimal(c *gin.Context) {
	UserID := c.GetUint("user_id")
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid animal ID"})
		return
	}
	animal, err := h.AnimalService.Get(UserID, uint(idUint))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Animal not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, animal)
}

func (h *AnimalHandler) UpdateAnimal(c *gin.Context) {
	UserID := c.GetUint("user_id")
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid animal ID"})
		return
	}
	var animal animals.Animal
	if err := c.ShouldBindJSON(&animal); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.AnimalService.Update(UserID, uint(idUint), &animal); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Animal not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	updated, _ := h.AnimalService.Get(UserID, uint(idUint))
	c.JSON(http.StatusOK, updated)
}

func (h *AnimalHandler) DeleteAnimal(c *gin.Context) {
	UserID := c.GetUint("user_id")
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid animal ID"})
		return
	}
	if err := h.AnimalService.Delete(UserID, uint(idUint)); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Animal not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Animal deleted successfully"})
}

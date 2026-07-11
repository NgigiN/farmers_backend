package animals

import (
	services "farm-backend/internal/services/animals"
	"farm-backend/internal/validation"
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
	var req validation.AnimalTypeRequest
	if err := validation.BindAndValidate(c, &req); err != nil {
		validation.RespondBindingError(c, err)
		return
	}
	entity := validation.AnimalTypeFromRequest(&req)
	if err := h.AnimalTypeService.Create(UserID, entity); err != nil {
		validation.RespondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, entity)
}

func (h *AnimalTypeHandler) ListAnimalTypes(c *gin.Context) {
	UserID := c.GetUint("user_id")
	animalTypes, err := h.AnimalTypeService.List(UserID)
	if err != nil {
		validation.RespondError(c, err)
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
			validation.RespondNotFound(c, "Animal type not found")
			return
		}
		validation.RespondError(c, err)
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
	var req validation.AnimalTypeRequest
	if err := validation.BindAndValidate(c, &req); err != nil {
		validation.RespondBindingError(c, err)
		return
	}
	entity := validation.AnimalTypeFromRequest(&req)
	if err := h.AnimalTypeService.Update(UserID, uint(idUint), entity); err != nil {
		if err == gorm.ErrRecordNotFound {
			validation.RespondNotFound(c, "Animal type not found")
			return
		}
		validation.RespondError(c, err)
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
			validation.RespondNotFound(c, "Animal type not found")
			return
		}
		validation.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Animal type deleted successfully"})
}
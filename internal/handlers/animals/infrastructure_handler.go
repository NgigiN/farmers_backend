package animals

import (
	services "farm-backend/internal/services/animals"
	"farm-backend/internal/validation"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type InfrastructureHandler struct {
	InfrastructureService *services.InfrastructureService
}

func NewInfrastructureHandler(infrastructureService *services.InfrastructureService) *InfrastructureHandler {
	return &InfrastructureHandler{InfrastructureService: infrastructureService}
}

func (h *InfrastructureHandler) CreateInfrastructure(c *gin.Context) {
	UserID := c.GetUint("user_id")
	var req validation.InfrastructureRequest
	if err := validation.BindAndValidate(c, &req); err != nil {
		validation.RespondBindingError(c, err)
		return
	}
	entity := validation.InfrastructureFromRequest(&req)
	if err := h.InfrastructureService.Create(UserID, entity); err != nil {
		validation.RespondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, entity)
}

func (h *InfrastructureHandler) ListInfrastructures(c *gin.Context) {
	UserID := c.GetUint("user_id")
	infrastructures, err := h.InfrastructureService.List(UserID)
	if err != nil {
		validation.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, infrastructures)
}

func (h *InfrastructureHandler) GetInfrastructure(c *gin.Context) {
	UserID := c.GetUint("user_id")
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid infrastructure ID"})
		return
	}
	infrastructure, err := h.InfrastructureService.Get(UserID, uint(idUint))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			validation.RespondNotFound(c, "Infrastructure not found")
			return
		}
		validation.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, infrastructure)
}

func (h *InfrastructureHandler) UpdateInfrastructure(c *gin.Context) {
	UserID := c.GetUint("user_id")
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid infrastructure ID"})
		return
	}
	var req validation.InfrastructureRequest
	if err := validation.BindAndValidate(c, &req); err != nil {
		validation.RespondBindingError(c, err)
		return
	}
	entity := validation.InfrastructureFromRequest(&req)
	if err := h.InfrastructureService.Update(UserID, uint(idUint), entity); err != nil {
		if err == gorm.ErrRecordNotFound {
			validation.RespondNotFound(c, "Infrastructure not found")
			return
		}
		validation.RespondError(c, err)
		return
	}
	updated, _ := h.InfrastructureService.Get(UserID, uint(idUint))
	c.JSON(http.StatusOK, updated)
}

func (h *InfrastructureHandler) DeleteInfrastructure(c *gin.Context) {
	UserID := c.GetUint("user_id")
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid infrastructure ID"})
		return
	}
	if err := h.InfrastructureService.Delete(UserID, uint(idUint)); err != nil {
		if err == gorm.ErrRecordNotFound {
			validation.RespondNotFound(c, "Infrastructure not found")
			return
		}
		validation.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Infrastructure deleted successfully"})
}
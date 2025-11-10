package animals

import (
	animals "farm-backend/internal/models/animals"
	services "farm-backend/internal/services/animals"
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
	var infrastructure animals.Infrastructure
	if err := c.ShouldBindJSON(&infrastructure); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.InfrastructureService.Create(UserID, &infrastructure); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, infrastructure)
}

func (h *InfrastructureHandler) ListInfrastructures(c *gin.Context) {
	UserID := c.GetUint("user_id")
	infrastructures, err := h.InfrastructureService.List(UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
			c.JSON(http.StatusNotFound, gin.H{"error": "Infrastructure not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
	var infrastructure animals.Infrastructure
	if err := c.ShouldBindJSON(&infrastructure); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.InfrastructureService.Update(UserID, uint(idUint), &infrastructure); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Infrastructure not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
			c.JSON(http.StatusNotFound, gin.H{"error": "Infrastructure not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Infrastructure deleted successfully"})
}

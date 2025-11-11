package plants

import (
	plantModels "farm-backend/internal/models/plants"
	plants "farm-backend/internal/services/plants"
	"strconv"

	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type LandHandler struct {
	LandService *plants.LandService
}

func NewLandHandler(landService *plants.LandService) *LandHandler {
	return &LandHandler{LandService: landService}
}

func (h *LandHandler) CreateLand(c *gin.Context) {
	UserID := c.GetUint("user_id")
	var land plantModels.Land
	if err := c.ShouldBindJSON(&land); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.LandService.Create(UserID, &land); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, land)
}

func (h *LandHandler) ListLands(c *gin.Context) {
	UserID := c.GetUint("user_id")
	lands, err := h.LandService.List(UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, lands)
}

func (h *LandHandler) GetLand(c *gin.Context) {
	UserID := c.GetUint("user_id")
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid land ID"})
		return
	}
	land, err := h.LandService.Get(UserID, uint(idUint))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Land not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, land)
}

func (h *LandHandler) UpdateLand(c *gin.Context) {
	UserID := c.GetUint("user_id")
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid land ID"})
		return
	}
	var land plantModels.Land
	if err := c.ShouldBindJSON(&land); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.LandService.Update(UserID, uint(idUint), &land); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Land not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	updated, _ := h.LandService.Get(UserID, uint(idUint))
	c.JSON(http.StatusOK, updated)
}

func (h *LandHandler) DeleteLand(c *gin.Context) {
	UserID := c.GetUint("user_id")
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid land ID"})
		return
	}
	if err := h.LandService.Delete(UserID, uint(idUint)); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Land not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Land deleted successfully"})
}

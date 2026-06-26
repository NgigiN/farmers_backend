package harvests

import (
	plantModels "farm-backend/internal/models/plants"
	plants "farm-backend/internal/services/plants"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HarvestHandler struct {
	HarvestService *plants.HarvestService
}

func NewHarvestHandler(harvestService *plants.HarvestService) *HarvestHandler {
	return &HarvestHandler{HarvestService: harvestService}
}

func (h *HarvestHandler) CreateHarvest(c *gin.Context) {
	userID := c.GetUint("user_id")
	var harvest plantModels.Harvest
	if err := c.ShouldBindJSON(&harvest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.HarvestService.Create(userID, &harvest); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, harvest)
}

// ListHarvests supports ?season_id=<id>
func (h *HarvestHandler) ListHarvests(c *gin.Context) {
	userID := c.GetUint("user_id")
	seasonIDStr := c.Query("season_id")

	var seasonID uint
	if seasonIDStr != "" {
		idUint, err := strconv.ParseUint(seasonIDStr, 10, 64)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid season_id"})
			return
		}
		seasonID = uint(idUint)
	}

	harvests, err := h.HarvestService.List(userID, seasonID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, harvests)
}

func (h *HarvestHandler) GetHarvest(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid harvest ID"})
		return
	}
	harvest, err := h.HarvestService.Get(userID, uint(idUint))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Harvest not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, harvest)
}

func (h *HarvestHandler) UpdateHarvest(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid harvest ID"})
		return
	}
	var harvest plantModels.Harvest
	if err := c.ShouldBindJSON(&harvest); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.HarvestService.Update(userID, uint(idUint), &harvest); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Harvest not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	updated, err := h.HarvestService.Get(userID, uint(idUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "harvest updated but could not be retrieved"})
		return
	}
	c.JSON(http.StatusOK, updated)
}

func (h *HarvestHandler) DeleteHarvest(c *gin.Context) {
	userID := c.GetUint("user_id")
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid harvest ID"})
		return
	}
	if err := h.HarvestService.Delete(userID, uint(idUint)); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Harvest not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Harvest deleted successfully"})
}
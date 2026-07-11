package harvests

import (
	plants "farm-backend/internal/services/plants"
	"farm-backend/internal/validation"
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
	var req validation.HarvestRequest
	if err := validation.BindAndValidate(c, &req); err != nil {
		validation.RespondBindingError(c, err)
		return
	}
	entity := validation.HarvestFromRequest(&req)
	if err := h.HarvestService.Create(userID, entity); err != nil {
		validation.RespondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, entity)
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
		validation.RespondError(c, err)
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
			validation.RespondNotFound(c, "Harvest not found")
			return
		}
		validation.RespondError(c, err)
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
	var req validation.HarvestRequest
	if err := validation.BindAndValidate(c, &req); err != nil {
		validation.RespondBindingError(c, err)
		return
	}
	entity := validation.HarvestFromRequest(&req)
	if err := h.HarvestService.Update(userID, uint(idUint), entity); err != nil {
		if err == gorm.ErrRecordNotFound {
			validation.RespondNotFound(c, "Harvest not found")
			return
		}
		validation.RespondError(c, err)
		return
	}
	updated, err := h.HarvestService.Get(userID, uint(idUint))
	if err != nil {
		validation.RespondError(c, err)
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
			validation.RespondNotFound(c, "Harvest not found")
			return
		}
		validation.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Harvest deleted successfully"})
}
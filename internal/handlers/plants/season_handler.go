package plants

import (
	plantModels "farm-backend/internal/models/plants"
	plants "farm-backend/internal/services/plants"
	"strconv"

	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type SeasonHandler struct {
	SeasonService *plants.SeasonService
}

func NewSeasonHandler(seasonService *plants.SeasonService) *SeasonHandler {
	return &SeasonHandler{SeasonService: seasonService}
}

func (h *SeasonHandler) CreateSeason(c *gin.Context) {
	UserID := c.GetUint("user_id")
	var season plantModels.Season
	if err := c.ShouldBindJSON(&season); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.SeasonService.Create(UserID, &season); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, season)
}

func (h *SeasonHandler) ListSeasons(c *gin.Context) {
	UserID := c.GetUint("user_id")
	seasons, err := h.SeasonService.List(UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, seasons)
}

func (h *SeasonHandler) GetSeason(c *gin.Context) {
	UserID := c.GetUint("user_id")
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid season ID"})
		return
	}
	season, err := h.SeasonService.Get(UserID, uint(idUint))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Season not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, season)
}

func (h *SeasonHandler) UpdateSeason(c *gin.Context) {
	UserID := c.GetUint("user_id")
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid season ID"})
		return
	}
	var season plantModels.Season
	if err := c.ShouldBindJSON(&season); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.SeasonService.Update(UserID, uint(idUint), &season); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Season not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	updated, _ := h.SeasonService.Get(UserID, uint(idUint))
	c.JSON(http.StatusOK, updated)
}

func (h *SeasonHandler) DeleteSeason(c *gin.Context) {
	UserID := c.GetUint("user_id")
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid season ID"})
		return
	}
	if err := h.SeasonService.Delete(UserID, uint(idUint)); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Season not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Season deleted successfully"})
}

package animals

import (
	animals "farm-backend/internal/models/animals"
	services "farm-backend/internal/services/animals"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type HerdHandler struct {
	HerdService *services.HerdService
}

func NewHerdHandler(herdService *services.HerdService) *HerdHandler {
	return &HerdHandler{HerdService: herdService}
}

func (h *HerdHandler) CreateHerd(c *gin.Context) {
	UserID := c.GetUint("user_id")
	var herd animals.Herd
	if err := c.ShouldBindJSON(&herd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.HerdService.Create(UserID, &herd); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, herd)
}

func (h *HerdHandler) ListHerds(c *gin.Context) {
	UserID := c.GetUint("user_id")
	herds, err := h.HerdService.List(UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, herds)
}

func (h *HerdHandler) GetHerd(c *gin.Context) {
	UserID := c.GetUint("user_id")
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid herd ID"})
		return
	}
	herd, err := h.HerdService.Get(UserID, uint(idUint))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Herd not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, herd)
}

func (h *HerdHandler) UpdateHerd(c *gin.Context) {
	UserID := c.GetUint("user_id")
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid herd ID"})
		return
	}
	var herd animals.Herd
	if err := c.ShouldBindJSON(&herd); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.HerdService.Update(UserID, uint(idUint), &herd); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Herd not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	updated, _ := h.HerdService.Get(UserID, uint(idUint))
	c.JSON(http.StatusOK, updated)
}

func (h *HerdHandler) DeleteHerd(c *gin.Context) {
	UserID := c.GetUint("user_id")
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid herd ID"})
		return
	}
	if err := h.HerdService.Delete(UserID, uint(idUint)); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Herd not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Herd deleted successfully"})
}

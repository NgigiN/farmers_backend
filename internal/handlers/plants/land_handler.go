package plants

import (
	plants "farm-backend/internal/services/plants"
	"farm-backend/internal/validation"
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
	var req validation.LandRequest
	if err := validation.BindAndValidate(c, &req); err != nil {
		validation.RespondBindingError(c, err)
		return
	}
	entity := validation.LandFromRequest(&req)
	if err := h.LandService.Create(UserID, entity); err != nil {
		validation.RespondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, entity)
}

func (h *LandHandler) ListLands(c *gin.Context) {
	UserID := c.GetUint("user_id")
	lands, err := h.LandService.List(UserID)
	if err != nil {
		validation.RespondError(c, err)
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
			validation.RespondNotFound(c, "Land not found")
			return
		}
		validation.RespondError(c, err)
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
	var req validation.LandRequest
	if err := validation.BindAndValidate(c, &req); err != nil {
		validation.RespondBindingError(c, err)
		return
	}
	entity := validation.LandFromRequest(&req)
	if err := h.LandService.Update(UserID, uint(idUint), entity); err != nil {
		if err == gorm.ErrRecordNotFound {
			validation.RespondNotFound(c, "Land not found")
			return
		}
		validation.RespondError(c, err)
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
			validation.RespondNotFound(c, "Land not found")
			return
		}
		validation.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Land deleted successfully"})
}
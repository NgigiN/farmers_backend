package animals

import (
	services "farm-backend/internal/services/animals"
	"farm-backend/internal/validation"
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
	var req validation.HerdRequest
	if err := validation.BindAndValidate(c, &req); err != nil {
		validation.RespondBindingError(c, err)
		return
	}
	entity := validation.HerdFromRequest(&req)
	if err := h.HerdService.Create(UserID, entity); err != nil {
		validation.RespondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, entity)
}

func (h *HerdHandler) ListHerds(c *gin.Context) {
	UserID := c.GetUint("user_id")
	herds, err := h.HerdService.List(UserID)
	if err != nil {
		validation.RespondError(c, err)
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
			validation.RespondNotFound(c, "Herd not found")
			return
		}
		validation.RespondError(c, err)
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
	var req validation.HerdRequest
	if err := validation.BindAndValidate(c, &req); err != nil {
		validation.RespondBindingError(c, err)
		return
	}
	entity := validation.HerdFromRequest(&req)
	if err := h.HerdService.Update(UserID, uint(idUint), entity); err != nil {
		if err == gorm.ErrRecordNotFound {
			validation.RespondNotFound(c, "Herd not found")
			return
		}
		validation.RespondError(c, err)
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
			validation.RespondNotFound(c, "Herd not found")
			return
		}
		validation.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Herd deleted successfully"})
}

func (h *HerdHandler) RecordActivity(c *gin.Context) {
	UserID := c.GetUint("user_id")
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid herd ID"})
		return
	}

	var req validation.HerdActivityRequest
	if err := validation.BindAndValidate(c, &req); err != nil {
		validation.RespondBindingError(c, err)
		return
	}
	entity := validation.HerdActivityFromRequest(&req)

	if err := h.HerdService.RecordActivity(UserID, uint(idUint), entity); err != nil {
		if err == gorm.ErrRecordNotFound {
			validation.RespondNotFound(c, "Herd not found")
			return
		}
		validation.RespondError(c, err)
		return
	}

	c.JSON(http.StatusCreated, entity)
}
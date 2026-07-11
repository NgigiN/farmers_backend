package activities

import (
	plants "farm-backend/internal/services/plants"
	"farm-backend/internal/validation"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ActivityHandler struct {
	ActivityService *plants.ActivityService
}

func NewActivityHandler(activityService *plants.ActivityService) *ActivityHandler {
	return &ActivityHandler{ActivityService: activityService}
}

func (h *ActivityHandler) CreateActivity(c *gin.Context) {
	UserID := c.GetUint("user_id")
	var req validation.ActivityRequest
	if err := validation.BindAndValidate(c, &req); err != nil {
		validation.RespondBindingError(c, err)
		return
	}
	entity := validation.ActivityFromRequest(&req)
	if err := h.ActivityService.Create(UserID, entity); err != nil {
		validation.RespondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, entity)
}

func (h *ActivityHandler) ListActivities(c *gin.Context) {
	UserID := c.GetUint("user_id")
	sourceType, err := validation.ValidateSourceType(c.Query("source_type"))
	if err != nil {
		validation.RespondError(c, err)
		return
	}
	activityList, err := h.ActivityService.List(UserID, sourceType)
	if err != nil {
		validation.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, activityList)
}

func (h *ActivityHandler) GetActivity(c *gin.Context) {
	UserID := c.GetUint("user_id")
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
		return
	}
	activity, err := h.ActivityService.Get(UserID, uint(idUint))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			validation.RespondNotFound(c, "Activity not found")
			return
		}
		validation.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, activity)
}

func (h *ActivityHandler) UpdateActivity(c *gin.Context) {
	UserID := c.GetUint("user_id")
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
		return
	}
	var req validation.ActivityRequest
	if err := validation.BindAndValidate(c, &req); err != nil {
		validation.RespondBindingError(c, err)
		return
	}
	entity := validation.ActivityFromRequest(&req)
	if err := h.ActivityService.Update(UserID, uint(idUint), entity); err != nil {
		if err == gorm.ErrRecordNotFound {
			validation.RespondNotFound(c, "Activity not found")
			return
		}
		validation.RespondError(c, err)
		return
	}
	updated, err := h.ActivityService.Get(UserID, uint(idUint))
	if err != nil {
		validation.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, updated)
}

func (h *ActivityHandler) DeleteActivity(c *gin.Context) {
	UserID := c.GetUint("user_id")
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid activity ID"})
		return
	}
	if err := h.ActivityService.Delete(UserID, uint(idUint)); err != nil {
		if err == gorm.ErrRecordNotFound {
			validation.RespondNotFound(c, "Activity not found")
			return
		}
		validation.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Activity deleted successfully"})
}
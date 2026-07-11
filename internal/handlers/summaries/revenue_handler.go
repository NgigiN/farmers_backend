package summaries

import (
	"errors"
	summaryModels "farm-backend/internal/models/summaries"
	summaries "farm-backend/internal/services/summaries"
	"farm-backend/internal/validation"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type RevenueHandler struct {
	RevenueService *summaries.RevenueService
}

func NewRevenueHandler(svc *summaries.RevenueService) *RevenueHandler {
	return &RevenueHandler{RevenueService: svc}
}

func (h *RevenueHandler) CreateRevenue(c *gin.Context) {
	UserID := c.GetUint("user_id")
	var req validation.RevenueRequest
	if err := validation.BindAndValidate(c, &req); err != nil {
		validation.RespondBindingError(c, err)
		return
	}
	entity := validation.RevenueFromRequest(&req)
	if err := h.RevenueService.Create(UserID, entity); err != nil {
		validation.RespondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, entity)
}

// ListRevenues supports:
//   - ?source=plant|animal
//   - ?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD
//
// Filtering is delegated to the service layer (SQL level).
func (h *RevenueHandler) ListRevenues(c *gin.Context) {
	UserID := c.GetUint("user_id")

	source, err := validation.ValidateRevenueSource(c.Query("source"))
	if err != nil {
		validation.RespondError(c, err)
		return
	}
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var revenues []summaryModels.Revenue
	var listErr error

	if startDateStr != "" && endDateStr != "" {
		startDate, err1 := time.Parse("2006-01-02", startDateStr)
		endDate, err2 := time.Parse("2006-01-02", endDateStr)
		if err1 != nil || err2 != nil {
			validation.RespondBindingError(c, errors.New("invalid date format — use YYYY-MM-DD"))
			return
		}
		revenues, listErr = h.RevenueService.ListByDateRange(UserID, startDate, endDate)
	} else if source != "" {
		revenues, listErr = h.RevenueService.ListBySource(UserID, source)
	} else {
		revenues, listErr = h.RevenueService.List(UserID)
	}

	if listErr != nil {
		validation.RespondError(c, listErr)
		return
	}
	c.JSON(http.StatusOK, revenues)
}

func (h *RevenueHandler) GetRevenue(c *gin.Context) {
	UserID := c.GetUint("user_id")
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		validation.RespondBindingError(c, errors.New("invalid revenue ID"))
		return
	}
	revenue, err := h.RevenueService.Get(UserID, uint(idUint))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			validation.RespondNotFound(c, "Revenue not found")
			return
		}
		validation.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, revenue)
}

func (h *RevenueHandler) UpdateRevenue(c *gin.Context) {
	UserID := c.GetUint("user_id")
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		validation.RespondBindingError(c, errors.New("invalid revenue ID"))
		return
	}
	var req validation.RevenueRequest
	if err := validation.BindAndValidate(c, &req); err != nil {
		validation.RespondBindingError(c, err)
		return
	}
	entity := validation.RevenueFromRequest(&req)
	if err := h.RevenueService.Update(UserID, uint(idUint), entity); err != nil {
		if err == gorm.ErrRecordNotFound {
			validation.RespondNotFound(c, "Revenue not found")
			return
		}
		validation.RespondError(c, err)
		return
	}
	updated, err := h.RevenueService.Get(UserID, uint(idUint))
	if err != nil {
		validation.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, updated)
}

func (h *RevenueHandler) DeleteRevenue(c *gin.Context) {
	UserID := c.GetUint("user_id")
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		validation.RespondBindingError(c, errors.New("invalid revenue ID"))
		return
	}
	if err := h.RevenueService.Delete(UserID, uint(idUint)); err != nil {
		if err == gorm.ErrRecordNotFound {
			validation.RespondNotFound(c, "Revenue not found")
			return
		}
		validation.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Revenue deleted successfully"})
}
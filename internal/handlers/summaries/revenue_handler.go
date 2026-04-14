package summaries

import (
	summaryModels "farm-backend/internal/models/summaries"
	summaries "farm-backend/internal/services/summaries"
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
	var revenue summaryModels.Revenue
	if err := c.ShouldBindJSON(&revenue); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.RevenueService.Create(UserID, &revenue); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, revenue)
}

// ListRevenues supports:
//   - ?source=plant|animal
//   - ?start_date=YYYY-MM-DD&end_date=YYYY-MM-DD
//
// Filtering is delegated to the service layer (SQL level).
func (h *RevenueHandler) ListRevenues(c *gin.Context) {
	UserID := c.GetUint("user_id")

	source := c.Query("source")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var revenues []summaryModels.Revenue
	var err error

	if startDateStr != "" && endDateStr != "" {
		startDate, err1 := time.Parse("2006-01-02", startDateStr)
		endDate, err2 := time.Parse("2006-01-02", endDateStr)
		if err1 != nil || err2 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format — use YYYY-MM-DD"})
			return
		}
		revenues, err = h.RevenueService.ListByDateRange(UserID, startDate, endDate)
	} else if source != "" {
		revenues, err = h.RevenueService.ListBySource(UserID, source)
	} else {
		revenues, err = h.RevenueService.List(UserID)
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, revenues)
}

func (h *RevenueHandler) GetRevenue(c *gin.Context) {
	UserID := c.GetUint("user_id")
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid revenue ID"})
		return
	}
	revenue, err := h.RevenueService.Get(UserID, uint(idUint))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Revenue not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, revenue)
}

func (h *RevenueHandler) UpdateRevenue(c *gin.Context) {
	UserID := c.GetUint("user_id")
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid revenue ID"})
		return
	}
	var revenue summaryModels.Revenue
	if err := c.ShouldBindJSON(&revenue); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.RevenueService.Update(UserID, uint(idUint), &revenue); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Revenue not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	updated, err := h.RevenueService.Get(UserID, uint(idUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "revenue updated but could not be retrieved"})
		return
	}
	c.JSON(http.StatusOK, updated)
}

func (h *RevenueHandler) DeleteRevenue(c *gin.Context) {
	UserID := c.GetUint("user_id")
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid revenue ID"})
		return
	}
	if err := h.RevenueService.Delete(UserID, uint(idUint)); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Revenue not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Revenue deleted successfully"})
}

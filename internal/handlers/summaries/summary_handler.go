package summaries

import (
	summaries "farm-backend/internal/models/summaries"
	services "farm-backend/internal/services/summaries"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CostCategoryHandler struct {
	CostCategoryService *services.CostCategoryService
}

func NewCostCategoryHandler(costCategoryService *services.CostCategoryService) *CostCategoryHandler {
	return &CostCategoryHandler{CostCategoryService: costCategoryService}
}

type RevenueHandler struct {
	RevenueService *services.RevenueService
}

func NewRevenueHandler(revenueService *services.RevenueService) *RevenueHandler {
	return &RevenueHandler{RevenueService: revenueService}
}

type AnalysisHandler struct {
	AnalysisService *services.AnalysisService
}

func NewAnalysisHandler(analysisService *services.AnalysisService) *AnalysisHandler {
	return &AnalysisHandler{AnalysisService: analysisService}
}

func (h *CostCategoryHandler) CreateCostCategory(c *gin.Context) {
	UserID := c.GetUint("user_id")
	var costCategory summaries.CostCategory
	if err := c.ShouldBindJSON(&costCategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.CostCategoryService.Create(UserID, &costCategory); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, costCategory)
}

func (h *CostCategoryHandler) ListCostCategories(c *gin.Context) {
	UserID := c.GetUint("user_id")
	costCategories, err := h.CostCategoryService.List(UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	typeFilter := c.Query("type")
	categoryFilter := c.Query("category")

	var filtered []summaries.CostCategory
	for _, cat := range costCategories {
		if typeFilter != "" && cat.Type != typeFilter {
			continue
		}
		if categoryFilter != "" && cat.Category != categoryFilter {
			continue
		}
		filtered = append(filtered, cat)
	}

	c.JSON(http.StatusOK, filtered)
}

func (h *CostCategoryHandler) GetCostCategory(c *gin.Context) {
	UserID := c.GetUint("user_id")
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cost category ID"})
		return
	}
	costCategory, err := h.CostCategoryService.Get(UserID, uint(idUint))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cost category not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, costCategory)
}

func (h *CostCategoryHandler) UpdateCostCategory(c *gin.Context) {
	UserID := c.GetUint("user_id")
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cost category ID"})
		return
	}
	var costCategory summaries.CostCategory
	if err := c.ShouldBindJSON(&costCategory); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.CostCategoryService.Update(UserID, uint(idUint), &costCategory); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cost category not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	updated, _ := h.CostCategoryService.Get(UserID, uint(idUint))
	c.JSON(http.StatusOK, updated)
}

func (h *CostCategoryHandler) DeleteCostCategory(c *gin.Context) {
	UserID := c.GetUint("user_id")
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid cost category ID"})
		return
	}
	if err := h.CostCategoryService.Delete(UserID, uint(idUint)); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Cost category not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Cost category deleted successfully"})
}

func (h *RevenueHandler) CreateRevenue(c *gin.Context) {
	UserID := c.GetUint("user_id")
	var revenue summaries.Revenue
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

func (h *RevenueHandler) ListRevenues(c *gin.Context) {
	UserID := c.GetUint("user_id")

	source := c.Query("source")
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	var revenues []summaries.Revenue
	var err error

	if startDateStr != "" && endDateStr != "" {
		startDate, err1 := time.Parse("2006-01-02", startDateStr)
		endDate, err2 := time.Parse("2006-01-02", endDateStr)
		if err1 != nil || err2 != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid date format. Use YYYY-MM-DD"})
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
	var revenue summaries.Revenue
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
	updated, _ := h.RevenueService.Get(UserID, uint(idUint))
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

func (h *AnalysisHandler) GetTotalCosts(c *gin.Context) {
	UserID := c.GetUint("user_id")
	total, err := h.AnalysisService.GetTotalCosts(UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total_costs": total})
}

func (h *AnalysisHandler) GetTotalRevenue(c *gin.Context) {
	UserID := c.GetUint("user_id")
	total, err := h.AnalysisService.GetTotalRevenue(UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"total_revenue": total})
}

func (h *AnalysisHandler) GetProfit(c *gin.Context) {
	UserID := c.GetUint("user_id")
	profit, err := h.AnalysisService.GetProfit(UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"profit": profit})
}

func (h *AnalysisHandler) GetCostBreakdown(c *gin.Context) {
	UserID := c.GetUint("user_id")
	breakdown, err := h.AnalysisService.GetCostBreakdownByCategory(UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, breakdown)
}

func (h *AnalysisHandler) GetRevenueBreakdown(c *gin.Context) {
	UserID := c.GetUint("user_id")
	breakdown, err := h.AnalysisService.GetRevenueBreakdownByType(UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, breakdown)
}

func (h *AnalysisHandler) GetMonthlySummary(c *gin.Context) {
	UserID := c.GetUint("user_id")
	yearStr := c.Query("year")
	if yearStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "year query parameter is required"})
		return
	}
	year, err := strconv.Atoi(yearStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid year format"})
		return
	}
	summary, err := h.AnalysisService.GetMonthlySummary(UserID, year)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, summary)
}

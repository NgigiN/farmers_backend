package summaries

import (
	summaries "farm-backend/internal/services/summaries"
	"farm-backend/internal/validation"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type AnalysisHandler struct {
	AnalysisService *summaries.AnalysisService
	CostService     *summaries.CostService
}

func NewAnalysisHandler(analysisSlice *summaries.AnalysisService, costSvc *summaries.CostService) *AnalysisHandler {
	return &AnalysisHandler{
		AnalysisService: analysisSlice,
		CostService:     costSvc,
	}
}

func (h *AnalysisHandler) GetTotalCosts(c *gin.Context) {
	UserID := c.GetUint("user_id")
	results, err := h.CostService.GetUnifiedDetailedCosts(UserID)
	if err != nil {
		validation.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, results)
}

func (h *AnalysisHandler) GetTotalRevenue(c *gin.Context) {
	UserID := c.GetUint("user_id")
	total, err := h.AnalysisService.GetTotalRevenue(UserID)
	if err != nil {
		validation.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"total_revenue": total})
}

func (h *AnalysisHandler) GetProfit(c *gin.Context) {
	UserID := c.GetUint("user_id")
	profit, err := h.AnalysisService.GetProfit(UserID)
	if err != nil {
		validation.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"profit": profit})
}

func (h *AnalysisHandler) GetCostBreakdown(c *gin.Context) {
	UserID := c.GetUint("user_id")
	breakdown, err := h.AnalysisService.GetCostBreakdownByCategory(UserID)
	if err != nil {
		validation.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, breakdown)
}

func (h *AnalysisHandler) GetRevenueBreakdown(c *gin.Context) {
	UserID := c.GetUint("user_id")
	breakdown, err := h.AnalysisService.GetRevenueBreakdownByType(UserID)
	if err != nil {
		validation.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, breakdown)
}

// GetMonthlySummary requires a ?year=YYYY query parameter.
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
		validation.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, summary)
}
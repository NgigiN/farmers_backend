package summaries

import (
	summaries "farm-backend/internal/services/summaries"
	"farm-backend/internal/validation"
	"net/http"

	"github.com/gin-gonic/gin"
)

type CostHandler struct {
	CostService *summaries.CostService
}

func NewCostHandler(svc *summaries.CostService) *CostHandler {
	return &CostHandler{CostService: svc}
}

func (h *CostHandler) GetTotalCostsBySeason(c *gin.Context) {
	UserID := c.GetUint("user_id")
	results, err := h.CostService.TotalCostBySeason(UserID)
	if err != nil {
		validation.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, results)
}
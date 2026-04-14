package summaries

import (
	summaryModels "farm-backend/internal/models/summaries"
	summaries "farm-backend/internal/services/summaries"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type CostCategoryHandler struct {
	CostCategoryService *summaries.CostCategoryService
}

func NewCostCategoryHandler(svc *summaries.CostCategoryService) *CostCategoryHandler {
	return &CostCategoryHandler{CostCategoryService: svc}
}

func (h *CostCategoryHandler) CreateCostCategory(c *gin.Context) {
	UserID := c.GetUint("user_id")
	var costCategory summaryModels.CostCategory
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

// ListCostCategories supports ?type=plant|animal&category=input|activity
// Filtering is pushed to the SQL layer via ListFiltered.
func (h *CostCategoryHandler) ListCostCategories(c *gin.Context) {
	UserID := c.GetUint("user_id")
	typeFilter := c.Query("type")
	categoryFilter := c.Query("category")

	costCategories, err := h.CostCategoryService.ListFiltered(UserID, typeFilter, categoryFilter)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, costCategories)
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
	var costCategory summaryModels.CostCategory
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
	updated, err := h.CostCategoryService.Get(UserID, uint(idUint))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "cost category updated but could not be retrieved"})
		return
	}
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

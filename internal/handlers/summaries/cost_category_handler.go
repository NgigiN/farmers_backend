package summaries

import (
	"errors"
	summaries "farm-backend/internal/services/summaries"
	"farm-backend/internal/validation"
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
	var req validation.CostCategoryRequest
	if err := validation.BindAndValidate(c, &req); err != nil {
		validation.RespondBindingError(c, err)
		return
	}
	entity := validation.CostCategoryFromRequest(&req)
	if err := h.CostCategoryService.Create(UserID, entity); err != nil {
		validation.RespondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, entity)
}

// ListCostCategories supports ?type=plant|animal&category=input|activity
// Filtering is pushed to the SQL layer via ListFiltered.
func (h *CostCategoryHandler) ListCostCategories(c *gin.Context) {
	UserID := c.GetUint("user_id")
	typeFilter, err := validation.ValidateCostCategoryType(c.Query("type"))
	if err != nil {
		validation.RespondError(c, err)
		return
	}
	categoryFilter, err := validation.ValidateCostCategoryCategory(c.Query("category"))
	if err != nil {
		validation.RespondError(c, err)
		return
	}

	costCategories, err := h.CostCategoryService.ListFiltered(UserID, typeFilter, categoryFilter)
	if err != nil {
		validation.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, costCategories)
}

func (h *CostCategoryHandler) GetCostCategory(c *gin.Context) {
	UserID := c.GetUint("user_id")
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		validation.RespondBindingError(c, errors.New("invalid cost category ID"))
		return
	}
	costCategory, err := h.CostCategoryService.Get(UserID, uint(idUint))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			validation.RespondNotFound(c, "Cost category not found")
			return
		}
		validation.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, costCategory)
}

func (h *CostCategoryHandler) UpdateCostCategory(c *gin.Context) {
	UserID := c.GetUint("user_id")
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		validation.RespondBindingError(c, errors.New("invalid cost category ID"))
		return
	}
	var req validation.CostCategoryRequest
	if err := validation.BindAndValidate(c, &req); err != nil {
		validation.RespondBindingError(c, err)
		return
	}
	entity := validation.CostCategoryFromRequest(&req)
	if err := h.CostCategoryService.Update(UserID, uint(idUint), entity); err != nil {
		if err == gorm.ErrRecordNotFound {
			validation.RespondNotFound(c, "Cost category not found")
			return
		}
		validation.RespondError(c, err)
		return
	}
	updated, err := h.CostCategoryService.Get(UserID, uint(idUint))
	if err != nil {
		validation.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, updated)
}

func (h *CostCategoryHandler) DeleteCostCategory(c *gin.Context) {
	UserID := c.GetUint("user_id")
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		validation.RespondBindingError(c, errors.New("invalid cost category ID"))
		return
	}
	if err := h.CostCategoryService.Delete(UserID, uint(idUint)); err != nil {
		if err == gorm.ErrRecordNotFound {
			validation.RespondNotFound(c, "Cost category not found")
			return
		}
		validation.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Cost category deleted successfully"})
}
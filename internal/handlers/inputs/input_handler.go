package inputs

import (
	plants "farm-backend/internal/services/plants"
	"farm-backend/internal/validation"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type InputHandler struct {
	InputService *plants.InputService
}

func NewInputHandler(inputService *plants.InputService) *InputHandler {
	return &InputHandler{InputService: inputService}
}

func (h *InputHandler) CreateInput(c *gin.Context) {
	UserID := c.GetUint("user_id")
	var req validation.InputRequest
	if err := validation.BindAndValidate(c, &req); err != nil {
		validation.RespondBindingError(c, err)
		return
	}
	entity := validation.InputFromRequest(&req)
	if err := h.InputService.Create(UserID, entity); err != nil {
		validation.RespondError(c, err)
		return
	}
	c.JSON(http.StatusCreated, entity)
}

func (h *InputHandler) ListInputs(c *gin.Context) {
	UserID := c.GetUint("user_id")
	sourceType, err := validation.ValidateSourceType(c.Query("source_type"))
	if err != nil {
		validation.RespondError(c, err)
		return
	}
	inputList, err := h.InputService.List(UserID, sourceType)
	if err != nil {
		validation.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, inputList)
}

func (h *InputHandler) GetInput(c *gin.Context) {
	UserID := c.GetUint("user_id")
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input ID"})
		return
	}
	input, err := h.InputService.Get(UserID, uint(idUint))
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			validation.RespondNotFound(c, "Input not found")
			return
		}
		validation.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, input)
}

func (h *InputHandler) UpdateInput(c *gin.Context) {
	UserID := c.GetUint("user_id")
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input ID"})
		return
	}
	var req validation.InputRequest
	if err := validation.BindAndValidate(c, &req); err != nil {
		validation.RespondBindingError(c, err)
		return
	}
	entity := validation.InputFromRequest(&req)
	if err := h.InputService.Update(UserID, uint(idUint), entity); err != nil {
		if err == gorm.ErrRecordNotFound {
			validation.RespondNotFound(c, "Input not found")
			return
		}
		validation.RespondError(c, err)
		return
	}
	updated, err := h.InputService.Get(UserID, uint(idUint))
	if err != nil {
		validation.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, updated)
}

func (h *InputHandler) DeleteInput(c *gin.Context) {
	UserID := c.GetUint("user_id")
	id := c.Param("id")
	idUint, err := strconv.ParseUint(id, 10, 64)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input ID"})
		return
	}
	if err := h.InputService.Delete(UserID, uint(idUint)); err != nil {
		if err == gorm.ErrRecordNotFound {
			validation.RespondNotFound(c, "Input not found")
			return
		}
		validation.RespondError(c, err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Input deleted successfully"})
}
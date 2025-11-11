package inputs

import (
	inputs "farm-backend/internal/models/plants"
	plants "farm-backend/internal/services/plants"
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
	var input inputs.Input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.InputService.Create(UserID, &input); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusCreated, input)
}

func (h *InputHandler) ListInputs(c *gin.Context) {
	UserID := c.GetUint("user_id")
	inputList, err := h.InputService.List(UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	sourceType := c.Query("source_type")
	if sourceType != "" {
		var filtered []inputs.Input
		for _, inp := range inputList {
			if inp.SourceType == sourceType {
				filtered = append(filtered, inp)
			}
		}
		c.JSON(http.StatusOK, filtered)
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
			c.JSON(http.StatusNotFound, gin.H{"error": "Input not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
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
	var input inputs.Input
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	if err := h.InputService.Update(UserID, uint(idUint), &input); err != nil {
		if err == gorm.ErrRecordNotFound {
			c.JSON(http.StatusNotFound, gin.H{"error": "Input not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	updated, _ := h.InputService.Get(UserID, uint(idUint))
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
			c.JSON(http.StatusNotFound, gin.H{"error": "Input not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Input deleted successfully"})
}

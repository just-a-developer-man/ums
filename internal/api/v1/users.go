package v1

import (
	"github.com/gin-gonic/gin"
)

// CreateUser handless the HTTP POST request to create a new user.
func (h *Handler) CreateUser(c *gin.Context) {
}

func (h *Handler) handleBody(c *gin.Context, dto any) error {
	err := c.ShouldBindJSON(dto)
	if err != nil {
		return &BindJSONError{err}
	}

	err = h.validator.ValidateStruct(dto)
	if err != nil {
	}
}

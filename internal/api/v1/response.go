package v1

import (
	"log/slog"
	"ums/internal/logger"

	"github.com/gin-gonic/gin"
)

// handleError is the wrapper to handle an error occured during request processing
// and to send appropriate response.
func (h *Handler) handleError(
	c *gin.Context,
	err error,
) {
	if err != nil {
		slog.ErrorContext(logger.ErrorCtx(c.Request.Context(), err), err.Error())
		errResp := mapError(err)
		c.JSON(errResp.Code, gin.H{"error": errResp.Desc})

	}
}

// handleOK is the wrapper to response with no error status.
func (h *Handler) handleOK(
	c *gin.Context,
	code int,
	response any,
) {
	if response != nil {
		c.JSON(code, response)
	} else {
		c.Status(code)
	}
}

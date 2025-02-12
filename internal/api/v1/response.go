package v1

import (
	"log/slog"
	"ums/internal/response"

	"github.com/gin-gonic/gin"
)

// handleError is the wrapper to handle an error occured during request processing
// and to send appropriate response.
func (h *Handler) handleError(
	c *gin.Context,
	logMessage string,
	err error,
) {
	if err != nil {
		slog.Error(logMessage, "err", err)
		errResp := response.MapError(err)
		c.JSON(errResp.Code, gin.H{"error": errResp.Desc})
	}
}

// handleOK is the wrapper to response with no error status.
func (h *Handler) handleOK(
	c *gin.Context,
	code int,
	response any,
	logMessage string,
) {
	slog.Info(logMessage)
	if response != nil {
		c.JSON(code, response)
	} else {
		c.Status(code)
	}
}

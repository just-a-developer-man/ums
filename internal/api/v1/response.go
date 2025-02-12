package v1

import (
	"log/slog"
	"ums/internal/response"

	"github.com/gin-gonic/gin"
)

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

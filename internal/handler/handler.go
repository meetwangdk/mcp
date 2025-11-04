package handler

import (
	"github.com/gin-gonic/gin"
	"github.com/google/wire"
	"hcs-agent/pkg/log"
)

var Set = wire.NewSet(
	NewHandler,
	NewDemoHandler,
)

type Handler struct {
	logger *log.Logger
}

func NewHandler(logger *log.Logger) *Handler {
	return &Handler{
		logger: logger,
	}
}
func GetUserIdFromCtx(ctx *gin.Context) string {
	_, exists := ctx.Get("claims")
	if !exists {
		return ""
	}
	return ""
}

package router

import (
	"hcs-agent/internal/handler"
	"hcs-agent/pkg/config"
	"hcs-agent/pkg/log"
)

type RouterDeps struct {
	Logger      *log.Logger
	Config      *config.Config
	DemoHandler *handler.DemoHandler
}

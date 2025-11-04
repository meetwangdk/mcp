package service

import (
	"github.com/google/wire"
	"hcs-agent/pkg/log"
	"hcs-agent/pkg/sid"
)

// ProviderSet is service providers.
var ProviderSet = wire.NewSet(NewService, NewDemoService)

type Service struct {
	logger *log.Logger
	sid    *sid.Sid
}

func NewService(
	logger *log.Logger,
	sid *sid.Sid,
) *Service {
	return &Service{
		logger: logger,
		sid:    sid,
	}
}

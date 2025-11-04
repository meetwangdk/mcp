package data

import (
	"context"
	"hcs-agent/internal/biz"
	"hcs-agent/internal/model"
	"hcs-agent/pkg/log"
)

var _ biz.DemoRepo = (*demoRepo)(nil)

type demoRepo struct {
	data *EsDb
	log  *log.Logger
}

// NewDemoRepo .
func NewDemoRepo(data *EsDb, logger *log.Logger) biz.DemoRepo {
	return &demoRepo{
		data: data,
		log:  logger,
	}
}

func (r *demoRepo) Create(context.Context, *model.Demo) (*model.Demo, error) {
	return &model.Demo{}, nil
}

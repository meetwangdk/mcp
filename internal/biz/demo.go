package biz

import (
	"context"
	"fmt"
	"hcs-agent/api/v1"
	"hcs-agent/internal/model"
	"hcs-agent/pkg/log"
)

// DemoRepo is a Greater repo.
type DemoRepo interface {
	Create(context.Context, *model.Demo) (*model.Demo, error)
}

// Demo is a Greeter.
type Demo struct {
	repo DemoRepo
	log  *log.Logger
}

// NewDemoUseCase new a Greeter.
func NewDemoUseCase(repo DemoRepo, logger *log.Logger) *Demo {
	return &Demo{repo: repo, log: logger}
}

// DemoCase creates a Greeter, and returns the new Greeter.
func (uc *Demo) DemoCase(ctx context.Context, s *v1.RegisterRequest) (*v1.Response, error) {
	uc.log.WithContext(ctx).Info(fmt.Sprintf("Get Demo email: %v", s.Email))
	// Transform the input into the domain entity.
	create, err := uc.repo.Create(ctx, &model.Demo{
		Email:    s.Email,
		Password: s.Password,
	})
	return &v1.Response{
		Data: create,
	}, err
}

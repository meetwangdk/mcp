package service

import (
	"context"
	v1 "hcs-agent/api/v1"
	"hcs-agent/internal/biz"
)

type DemoService interface {
	Register(ctx context.Context, req *v1.RegisterRequest) error
	Login(ctx context.Context, req *v1.LoginRequest) (string, error)
	GetProfile(ctx context.Context, userId string) (*v1.GetProfileResponseData, error)
	UpdateProfile(ctx context.Context, userId string, req *v1.UpdateProfileRequest) error
}

func NewDemoService(service *Service, demoBiz *biz.Demo) DemoService {
	return &demoService{
		demoBiz: demoBiz,
		Service: service,
	}
}

type demoService struct {
	*Service
	demoBiz *biz.Demo
}

func (s *demoService) Register(ctx context.Context, req *v1.RegisterRequest) error {
	s.demoBiz.DemoCase(ctx, req)
	return nil
}

func (s *demoService) Login(ctx context.Context, req *v1.LoginRequest) (string, error) {

	return "token", nil
}

func (s *demoService) GetProfile(ctx context.Context, userId string) (*v1.GetProfileResponseData, error) {

	return &v1.GetProfileResponseData{
		UserId:   "user.UserId",
		Nickname: "user.Nickname",
	}, nil
}

func (s *demoService) UpdateProfile(ctx context.Context, userId string, req *v1.UpdateProfileRequest) error {
	return nil
}

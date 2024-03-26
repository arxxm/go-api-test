package service

import (
	"context"
	"go-api-test/pkg/domain"
)

type Service struct {
	Users Users
}

func NewService(users Users) *Service {
	return &Service{
		Users: users,
	}
}

type Users interface {
	Create(ctx context.Context, user domain.User) (int64, error)
	Update(ctx context.Context, id int64, user domain.User) error
	GetList(ctx context.Context, params *domain.UsersParam) ([]domain.User, uint64, error)
	Delete(ctx context.Context, id int64) error
	GetByID(ctx context.Context, id int64) (domain.User, error)
}

package redis

import (
	"context"
	"go-api-test/pkg/domain"
)

type Repository struct {
	Users Users
}

func NewRepository(users *UsersStorage) *Repository {
	return &Repository{
		Users: users,
	}
}

type Users interface {
	GetByQuery(ctx context.Context, params *domain.UsersParam) ([]domain.User, uint64, error)
	Create(ctx context.Context, params *domain.UsersParam, users []domain.User) error
}

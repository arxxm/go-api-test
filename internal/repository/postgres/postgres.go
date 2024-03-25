package postgres

import (
	"context"
	"go-api-test/pkg/domain"
)

type Repository struct {
	Users Users
}

func NewRepository(users Users) *Repository {
	return &Repository{
		Users: users,
	}
}

type Users interface {
	GetList(ctx context.Context, params *domain.UsersParam) ([]domain.User, uint64, error)
	GetByID(ctx context.Context, id int64) (domain.User, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, id int64, user domain.User) error
	Create(ctx context.Context, user domain.User) (int64, error)
}

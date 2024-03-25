package service

import (
	"context"
	"go-api-test/internal/repository"
	"go-api-test/pkg/domain"
	"log"
)

type UsersService struct {
	repo *repository.Repository
}

func NewUsersService(repo *repository.Repository) *UsersService {
	return &UsersService{
		repo: repo,
	}
}

func (s *UsersService) Create(ctx context.Context, user domain.User) (int64, error) {
	id, err := s.repo.Postgres.Users.Create(ctx, user)
	if err != nil {
		log.Printf("[ERROR] User create err: %v id: %d\n", err, id)
		return 0, err
	}
	return id, nil
}

func (s *UsersService) Update(ctx context.Context, id int64, person domain.User) error {
	if err := s.repo.Postgres.Users.Update(ctx, id, person); err != nil {
		log.Printf("[ERROR] User update err: %v, id: %d\n", err, id)
		return err
	}
	return nil
}

func (s *UsersService) Delete(ctx context.Context, id int64) error {
	if err := s.repo.Postgres.Users.Delete(ctx, id); err != nil {
		log.Printf("[ERROR] User delete err: %v, id: %d\n", err, id)
		return err
	}
	return nil
}

func (s *UsersService) GetByID(ctx context.Context, id int64) (domain.User, error) {
	var user domain.User
	var err error

	if user, err = s.repo.Postgres.Users.GetByID(ctx, id); err != nil {
		log.Printf("[ERROR] User get err: %v, id: %d\n", err, id)
		return user, err
	}
	return user, nil
}

func (s *UsersService) GetList(ctx context.Context, params *domain.UsersParam) ([]domain.User, uint64, error) {
	persons, total, err := s.repo.Postgres.Users.GetList(ctx, params)
	if err != nil {
		log.Printf("[ERROR] User get list err: %v\n", err)
		return nil, 0, err
	}
	return persons, total, err
}

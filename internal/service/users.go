package service

import (
	"context"
	"go-api-test/internal/repository"
	"go-api-test/internal/service/cache"
	"go-api-test/pkg/domain"
	"log"
)

type UsersService struct {
	repo       *repository.Repository
	usersCache *cache.UsersCache
}

func NewUsersService(repo *repository.Repository) *UsersService {
	return &UsersService{
		repo:       repo,
		usersCache: cache.NewUsersCache(),
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

func (s *UsersService) Update(ctx context.Context, id int64, user domain.User) error {
	if err := s.repo.Postgres.Users.Update(ctx, id, user); err != nil {
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
	s.usersCache.Delete(id)
	return nil
}

func (s *UsersService) GetByID(ctx context.Context, id int64) (domain.User, error) {
	var user domain.User
	var err error

	user, ok := s.usersCache.Get(id)
	if ok {
		return user, nil
	}

	if user, err = s.repo.Postgres.Users.GetByID(ctx, id); err != nil {
		log.Printf("[ERROR] User get err: %v, id: %d\n", err, id)
		return user, err
	}

	s.usersCache.Set(id, user)

	return user, nil
}

func (s *UsersService) GetList(ctx context.Context, params *domain.UsersParam) ([]domain.User, uint64, error) {

	var users []domain.User
	var total uint64
	var err error

	users, total, err = s.repo.Redis.Users.GetByQuery(ctx, params)
	if err != nil {
		log.Printf("[INFO] User get list from cache err: %v\n", err)
	}
	if total > 0 {
		return users, total, nil
	}

	users, total, err = s.repo.Postgres.Users.GetList(ctx, params)
	if err != nil {
		log.Printf("[ERROR] User get list err: %v\n", err)
		return nil, 0, err
	}

	if err := s.repo.Redis.Users.Create(ctx, params, users); err != nil {
		log.Printf("[ERROR] User set cache list err: %v\n", err)
	}

	return users, total, err
}

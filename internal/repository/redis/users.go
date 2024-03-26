package redis

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/go-redis/redis/v8"
	"go-api-test/pkg/domain"
	"time"
)

type UsersStorage struct {
	conn *redis.Client
}

func NewUsersStorage(conn *redis.Client) *UsersStorage {
	return &UsersStorage{conn: conn}
}

func (r *UsersStorage) Create(ctx context.Context, params *domain.UsersParam, users []domain.User) error {
	key := fmt.Sprintf("users:%+v", params)
	valBytes, err := json.Marshal(users)
	if err != nil {
		return fmt.Errorf("failed to marshal users: %v", err)
	}

	if err := r.conn.Set(ctx, key, valBytes, 5*time.Minute).Err(); err != nil {
		return fmt.Errorf("failed to set value to cache: %v", err)
	}

	return nil
}

func (r *UsersStorage) GetByQuery(ctx context.Context, params *domain.UsersParam) ([]domain.User, uint64, error) {
	key := fmt.Sprintf("users:%+v", params)
	val, err := r.conn.Get(ctx, key).Bytes()
	if err == redis.Nil {
		return nil, 0, fmt.Errorf("users not found in cache")
	} else if err != nil {
		return nil, 0, fmt.Errorf("failed to get value from cache: %v", err)
	}

	var users []domain.User
	if err := json.Unmarshal(val, &users); err != nil {
		return nil, 0, fmt.Errorf("failed to unmarshal users: %v", err)
	}

	return users, uint64(len(users)), nil
}

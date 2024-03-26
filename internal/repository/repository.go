package repository

import (
	postgresrepo "go-api-test/internal/repository/postgres"
	redisrepo "go-api-test/internal/repository/redis"
)

type Repository struct {
	Postgres *postgresrepo.Repository
	Redis    *redisrepo.Repository
}

func NewRepository(postgres *postgresrepo.Repository, redis *redisrepo.Repository) *Repository {
	return &Repository{
		Postgres: postgres,
		Redis:    redis,
	}
}

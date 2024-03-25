package repository

import (
	postgresrepo "go-api-test/internal/repository/postgres"
)

type Repository struct {
	Postgres *postgresrepo.Repository
}

func NewRepository(postgres *postgresrepo.Repository) *Repository {
	return &Repository{
		Postgres: postgres,
	}
}

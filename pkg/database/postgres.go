package database

import (
	"context"
	"fmt"
	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
	"go-api-test/config"
)

type Postgres struct {
	cfg  *config.Config
	conn *pgx.Conn
}

func NewPostgres(cfg *config.Config) *Postgres {
	return &Postgres{cfg: cfg}
}

func (p *Postgres) Init() (*pgx.Conn, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable",
		p.cfg.PostgresUser, p.cfg.PostgresPass, p.cfg.PostgresHost, p.cfg.PostgresName)

	conn, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

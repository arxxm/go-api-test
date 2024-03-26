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

	_, err = conn.Exec(context.Background(), `
        CREATE TABLE IF NOT EXISTS users (
            id SERIAL PRIMARY KEY,
            name VARCHAR(255) NOT NULL,
            last_name VARCHAR(255) NOT NULL,
            surname VARCHAR(255),
            gender VARCHAR(20) NOT NULL,
            status VARCHAR(20) NOT NULL,
            date_of_birth TIMESTAMP,
            created_at TIMESTAMP NOT NULL
        );
    `)
	if err != nil {
		return nil, err
	}

	_, err = conn.Exec(context.Background(), `
        CREATE INDEX IF NOT EXISTS idx_user_id ON users (id);
        CREATE INDEX IF NOT EXISTS idx_user_gender ON users (gender);
        CREATE INDEX IF NOT EXISTS idx_user_status ON users (status);
    `)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

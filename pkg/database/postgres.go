package database

import (
	"database/sql"
	"fmt"
	_ "github.com/lib/pq"
	"go-api-test/config"
)

type Postgres struct {
	cfg  *config.Config
	conn *sql.DB
}

func NewPostgres(cfg *config.Config) *Postgres {
	return &Postgres{cfg: cfg}
}

func (p *Postgres) Init() (*sql.DB, error) {
	addr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", p.cfg.PostgresUser, p.cfg.PostgresPass, p.cfg.PostgresHost, p.cfg.PostgresName)

	db, err := sql.Open("postgres", addr)
	if err != nil {
		return nil, err
	}

	db.SetMaxOpenConns(200)
	db.SetMaxIdleConns(10)

	err = db.Ping()
	if err != nil {
		return nil, err
	}

	return db, nil
}

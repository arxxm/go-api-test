package main

import (
	"context"
	"database/sql"
	_ "github.com/lib/pq"
	"go-api-test/config"
	"go-api-test/internal/repository"
	"go-api-test/internal/repository/postgres"
	"go-api-test/internal/server"
	"go-api-test/internal/service"
	"go-api-test/internal/transport/rest"
	"go-api-test/pkg/database"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {

	var cfg *config.Config
	var err error
	if cfg, err = config.New(); err != nil {
		log.Fatal("config err:", err)
	}

	var (
		_, cancel = context.WithCancel(context.Background())
		quit      = make(chan os.Signal, 1)
	)

	p := database.NewPostgres(cfg)
	pConn, err := p.Init()
	if err != nil {
		log.Fatal(err)
	}

	defer func(pConn *sql.DB) {
		if err = pConn.Close(); err != nil {
			log.Fatal(err)
		}
	}(pConn)

	pqUsersRepo := postgres.NewUsersStorage(pConn)
	pqRepo := postgres.NewRepository(pqUsersRepo)
	repo := repository.NewRepository(pqRepo)

	usersService := service.NewUsersService(repo)
	s := service.NewService(usersService)
	h := rest.NewHandler(s)
	engine := h.InitRoutes()

	srv := server.New(8080, engine)

	go func() {
		if err := srv.Run(); err != nil {
			log.Fatal(err)
		}
	}()

	log.Println("service started")

	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)

	sig := <-quit

	ctxSrvStop, cancelSrvStop := context.WithTimeout(context.Background(), 5*time.Second)

	if err := srv.Stop(ctxSrvStop); err != nil {
		log.Fatal(err)
	}

	cancelSrvStop()
	cancel()

	log.Printf("admin shutdown signal %v", sig)
}

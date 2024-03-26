package main

import (
	"context"
	rc "github.com/go-redis/redis/v8"
	"github.com/jackc/pgx/v5"
	_ "github.com/lib/pq"
	"go-api-test/config"
	"go-api-test/internal/repository"
	"go-api-test/internal/repository/postgres"
	"go-api-test/internal/repository/redis"
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
		ctx, cancel = context.WithCancel(context.Background())
		quit        = make(chan os.Signal, 1)
	)

	p := database.NewPostgres(cfg)
	pConn, err := p.Init()
	if err != nil {
		log.Fatal(err)
	}

	defer func(pConn *pgx.Conn) {
		if err = pConn.Close(ctx); err != nil {
			log.Fatal(err)
		}
	}(pConn)

	r := database.NewRedis(cfg)
	rConn, err := r.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer func(pConn *rc.Client) {
		if err = pConn.Close(); err != nil {
			log.Fatal(err)
		}
	}(rConn)

	pqUsersRepo := postgres.NewUsersStorage(pConn)
	redisUsersRepo := redis.NewUsersStorage(rConn)
	pqRepo := postgres.NewRepository(pqUsersRepo)
	redisRepo := redis.NewRepository(redisUsersRepo)
	repo := repository.NewRepository(pqRepo, redisRepo)

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

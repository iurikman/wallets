package main

import (
	"context"
	"os/signal"
	"syscall"

	"github.com/iurikman/wallets/internal/config"
	"github.com/iurikman/wallets/internal/rest"
	"github.com/iurikman/wallets/internal/service"
	"github.com/iurikman/wallets/internal/store"
	migrate "github.com/rubenv/sql-migrate"
	log "github.com/sirupsen/logrus"
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGQUIT, syscall.SIGHUP)
	defer cancel()

	cfg := config.NewConfig()

	db, err := store.New(ctx, store.Config{
		PGUser:     cfg.PostgresUser,
		PGPass:     cfg.PostgresPassword,
		PGHost:     cfg.PostgresHost,
		PGPort:     cfg.PostgresPort,
		PGDatabase: cfg.PostgresDatabase,
	})
	if err != nil {
		log.Panicf("store.NewPostgres(context.Background(), store.ServerConfig{...} err: %v", err)
	}

	if err := db.Migrate(migrate.Up); err != nil {
		log.Panicf("pgStore.Migrate: %v", err)
	}

	log.Info("successful migration")

	svc := service.New(db)

	srv, err := rest.NewServer(
		rest.ServerConfig{BindAddress: cfg.BindAddress},
		svc,
	)
	if err != nil {
		log.Panicf("rest.NewServer(cfg) err: %v", err)
	}

	if err := srv.Start(ctx); err != nil {
		log.Panicf("srv.Start(ctx) err: %v", err)
	}

	log.Info("service stopped")
}

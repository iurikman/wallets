package rest

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/go-chi/chi"
	"github.com/sirupsen/logrus"
)

type ServerConfig struct {
	BindAddress string
}

const (
	readHeaderTimeout       = 10 * time.Second
	maxHeaderBytes          = 1 << 20
	gracefulShutdownTimeout = 5 * time.Second
)

type Server struct {
	serverConfig ServerConfig
	service      service
	router       *chi.Mux
	server       *http.Server
}

func NewServer(serverConfig ServerConfig, srv service) (*Server, error) {
	router := chi.NewRouter()

	return &Server{
		serverConfig: serverConfig,
		service:      srv,
		router:       router,
		server: &http.Server{
			Addr:              serverConfig.BindAddress,
			Handler:           router,
			ReadHeaderTimeout: readHeaderTimeout,
			MaxHeaderBytes:    maxHeaderBytes,
		},
	}, nil
}

func (s *Server) Start(ctx context.Context) error {
	s.configRouter()

	go func() {
		<-ctx.Done()
		ctxWithTimeout, cancel := context.WithTimeout(ctx, gracefulShutdownTimeout)

		defer cancel()

		err := s.server.Shutdown(ctxWithTimeout)
		if err != nil {
			logrus.Warnf("failed to shutdown gracefully %s", err)
		}
	}()

	err := s.server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		return fmt.Errorf("s.server.ListenAndServe() err: %w", err)
	}

	return nil
}

func (s *Server) configRouter() {
	s.router.Route("/api/v1", func(r chi.Router) {
		r.Route("/wallets", func(r chi.Router) {
			r.Post("/", s.createWallet)
			r.Get("/{id}", s.getWallet)

			r.Put("/withdraw", s.withdraw)
			r.Put("/deposit", s.deposit)
		})
	})
}

package server

import (
	"context"
	"log"
	"net"
	"net/http"
	"time"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/julienschmidt/httprouter"
	"github.com/s1ovac/order-subscribe/internal/store/config"
	"github.com/sirupsen/logrus"
)

type Server struct {
	ctx    context.Context
	router *httprouter.Router
	logger *logrus.Logger
	config *config.ServerConfig
	db     *pgxpool.Pool
	srv    *http.Server
}

func NewServer(ctx context.Context, router *httprouter.Router, logger *logrus.Logger, config *config.ServerConfig, db *pgxpool.Pool) *Server {
	return &Server{
		ctx:    ctx,
		router: router,
		logger: logger,
		config: config,
		db:     db,
		srv:    &http.Server{},
	}
}

func (s *Server) Start() {
	s.logger.Info("starting server")
	listener, err := net.Listen(s.config.Protocol, s.config.BindAddress)
	if err != nil {
		panic(err)
	}
	s.srv = &http.Server{
		Handler:      s.router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	s.logger.Fatal(s.srv.Serve(listener))
}

func (s *Server) Shutdown() {
	s.logger.Info("server stopped")

	ctxShutdown, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer func() {
		cancel()
	}()

	var err error
	if err = s.srv.Shutdown(ctxShutdown); err != nil {
		log.Fatalf("server Shutdown Failed %q", err)
	}

	s.logger.Printf("server exited properly")

	if err == http.ErrServerClosed {
		err = nil
	}
}

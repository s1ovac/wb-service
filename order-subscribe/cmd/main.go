package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/s1ovac/order-subscribe/internal/cache"
	"github.com/s1ovac/order-subscribe/internal/pkg/logging"
	"github.com/s1ovac/order-subscribe/internal/store/config"
	"github.com/s1ovac/order-subscribe/internal/store/databases/order"
	"github.com/s1ovac/order-subscribe/internal/store/databases/order/handler"
	"github.com/s1ovac/order-subscribe/internal/store/databases/postgresql"
	"github.com/s1ovac/order-subscribe/internal/subscribe"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logging.Init()
	router := httprouter.New()
	logger.Info("create router")
	cfgDataBase := config.NewStorageConfig()
	cfgServer := config.NewServerConfig()
	ctx, cancel := context.WithCancel(context.Background())
	postgreSQL, err := postgresql.NewClient(ctx, 3, cfgDataBase)
	if err != nil {
		logger.Fatal(err)
	}
	rep := order.NewRepository(postgreSQL)
	sb := subscribe.New(logger, rep, postgreSQL)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		oscall := <-c
		logger.Info("system call:%+v", oscall)
		cancel()
	}()

	ch := cache.NewCache(rep)
	if err := ch.InitCache(ctx); err != nil {
		logger.Fatal(err)
	}
	orderHandler := handler.NewHandler(ctx, &rep, logger, ch, sb)
	orderHandler.Register(ctx, router)
	start(router, logger, cfgServer)
}

func start(router *httprouter.Router, logger *logrus.Logger, config *config.ServerConfig) {
	listener, err := net.Listen(config.Protocol, config.BindAddress)
	if err != nil {
		panic(err)
	}
	server := http.Server{
		Handler:      router,
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	logger.Fatal(server.Serve(listener))
}

package main

import (
	"context"
	"net"
	"net/http"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/s1ovac/order-subscribe/internal/cache"
	"github.com/s1ovac/order-subscribe/internal/pkg/logging"
	"github.com/s1ovac/order-subscribe/internal/store/config"
	"github.com/s1ovac/order-subscribe/internal/store/databases/order"
	"github.com/s1ovac/order-subscribe/internal/store/databases/order/handler"
	"github.com/s1ovac/order-subscribe/internal/store/databases/postgresql"
	"github.com/sirupsen/logrus"
)

func main() {
	logger := logging.Init()
	router := httprouter.New()
	logger.Info("create router")
	// sb := subscribe.New(logger)
	// newOrder, err := sb.SubscribeToChannel()
	// if err != nil {
	// 	logger.Fatal(err)
	// }
	cfgDataBase := config.NewStorageConfig()
	cfgServer := config.NewServerConfig()
	postgreSQL, err := postgresql.NewClient(context.TODO(), 3, cfgDataBase)
	if err != nil {
		logger.Fatal(err)
	}
	rep := order.NewRepository(postgreSQL)
	ch := cache.NewCache(rep)
	if err := ch.InitCache(context.TODO()); err != nil {
		logger.Fatal(err)
	}
	orderHandler := handler.NewHandler(&rep, logger, ch)
	orderHandler.Register(router)
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

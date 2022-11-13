package main

import (
	"context"
	"net"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/julienschmidt/httprouter"
	"github.com/nats-io/stan.go"
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
	ch := cache.NewCache(rep)
	if err := ch.InitCache(ctx); err != nil {
		logger.Fatal(err)
	}
	sb := subscribe.New(logger, rep, postgreSQL, ctx)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		oscall := <-c
		logger.Info("system call:%+v", oscall)
		cancel()
	}()
	conn, err := stan.Connect(sb.ClusterID, sb.ClientID, stan.NatsURL("nats://localhost:4222"))
	if err != nil {
		logger.Fatal(err)
	}
	defer conn.Close()

	if _, err = conn.Subscribe(sb.Channel, sb.CreateOrder, stan.StartWithLastReceived()); err != nil {
		return
	}
	orderHandler := handler.NewHandler(ctx, &rep, logger, ch, sb)
	orderHandler.Register(ctx, router)
	start(ctx, router, logger, cfgServer)
}

func start(ctx context.Context, router *httprouter.Router, logger *logrus.Logger, config *config.ServerConfig) {
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

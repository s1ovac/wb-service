package main

import (
	"context"
	"os"
	"os/signal"

	"github.com/julienschmidt/httprouter"
	"github.com/nats-io/stan.go"
	"github.com/s1ovac/order-subscribe/internal/cache"
	"github.com/s1ovac/order-subscribe/internal/pkg/logging"
	"github.com/s1ovac/order-subscribe/internal/server"
	"github.com/s1ovac/order-subscribe/internal/store/config"
	"github.com/s1ovac/order-subscribe/internal/store/databases/order"
	"github.com/s1ovac/order-subscribe/internal/store/databases/order/handler"
	"github.com/s1ovac/order-subscribe/internal/store/databases/postgresql"
	"github.com/s1ovac/order-subscribe/internal/subscribe"
)

func main() {
	logger := logging.Init()
	router := httprouter.New()
	logger.Info("create router")
	cfgDataBase := config.NewStorageConfig()
	cfgServer := config.NewServerConfig()
	ctx, cancel := context.WithCancel(context.Background())
	db, err := postgresql.NewClient(ctx, 3, cfgDataBase)
	if err != nil {
		logger.Fatal(err)
	}
	rep := order.NewRepository(db)
	ch := cache.NewCache(rep)
	if err := ch.InitCache(ctx); err != nil {
		logger.Fatal(err)
	}
	server := server.NewServer(ctx, router, logger, cfgServer, db)
	sb := subscribe.New(logger, rep, db, ctx)
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		oscall := <-c
		logger.Info("system call:%+v", oscall)
		server.Shutdown()
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

	server.Start()
}

// func start(ctx context.Context, router *httprouter.Router, logger *logrus.Logger, config *config.ServerConfig) {
// 	listener, err := net.Listen(config.Protocol, config.BindAddress)
// 	if err != nil {
// 		panic(err)
// 	}
// 	server := http.Server{
// 		Handler:      router,
// 		WriteTimeout: 15 * time.Second,
// 		ReadTimeout:  15 * time.Second,
// 	}
// 	logger.Fatal(server.Serve(listener))
// }

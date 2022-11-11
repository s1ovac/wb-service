package main

import (
	"context"
	"fmt"
	"log"

	"github.com/julienschmidt/httprouter"
	"github.com/s1ovac/order-subscribe/internal/cache"
	"github.com/s1ovac/order-subscribe/internal/store/config"
	"github.com/s1ovac/order-subscribe/internal/store/databases/order"
	"github.com/s1ovac/order-subscribe/internal/store/databases/postgresql"
	"github.com/s1ovac/order-subscribe/internal/subscribe"
)

func main() {
	router := httprouter.New()
	sb := subscribe.New()
	newOrder, err := sb.SubscribeToChannel()
	if err != nil {
		log.Fatal(err)
	}

	cfg := config.NewConfig()
	postgreSQL, err := postgresql.NewClient(context.TODO(), 3, cfg)
	if err != nil {
		log.Fatal(err)
	}
	rep := order.NewRepository(postgreSQL)
	err = rep.Create(context.TODO(), newOrder, postgreSQL)
	if err != nil {
		log.Fatal(err)
	}
	orders, err := rep.FindAll(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	c := cache.NewCache(newOrder)
	err = c.InitCache(newOrder.OrderUID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(orders)
}

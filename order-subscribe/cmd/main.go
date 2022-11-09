package main

import (
	"context"
	"fmt"
	"log"

	"github.com/s1ovac/order-subscribe/internal/cache"
	"github.com/s1ovac/order-subscribe/internal/store/config"
	"github.com/s1ovac/order-subscribe/internal/store/databases/order"
	"github.com/s1ovac/order-subscribe/internal/store/databases/postgresql"
	//"github.com/s1ovac/order-subscribe/internal/subscribe"
)

func main() {
	// sb := subscribe.New()
	// newOrder, err := sb.SubscribeToChannel()
	// if err != nil {
	// 	log.Fatal(err)
	// }
	cfg := config.NewConfig()
	postgreSQL, err := postgresql.NewClient(context.TODO(), 3, cfg)
	if err != nil {
		log.Fatal(err)
	}
	rep := order.NewRepository(postgreSQL)
	ord, err := rep.FindOne(context.TODO(), "74e72147-1b20-4cee-8056-4868acfe639a")
	if err != nil {
		log.Fatal(err)
	}
	c := cache.NewCache(&ord)
	err = c.InitCache(ord.OrderUID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(c.Cache(ord.OrderUID))
	fmt.Println(ord)
}

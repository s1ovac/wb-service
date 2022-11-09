package main

import (
	"context"
	"fmt"
	"log"

	"github.com/s1ovac/order-subscribe/internal/store/config"
	"github.com/s1ovac/order-subscribe/internal/store/databases/order"
	"github.com/s1ovac/order-subscribe/internal/store/databases/postgresql"
	//"github.com/s1ovac/order-subscribe/internal/subscribe"
)

func main() {
	//sb := subscribe.New()
	//newOrder, err := sb.SubscribeToChannel()
	//if err != nil {
	//	log.Fatal(err)
	//}
	cfg := config.NewConfig()
	postgreSQL, err := postgresql.NewClient(context.TODO(), 3, cfg)
	if err != nil {
		log.Fatal(err)
	}
	rep := order.NewRepository(postgreSQL)
	orders, err := rep.FindAll(context.TODO())
	if err != nil {
		log.Fatal(err)
	}
	for _, order := range orders {
		fmt.Println(order)
	}
}

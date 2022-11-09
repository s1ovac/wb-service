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
	order, err := rep.FindOne(context.TODO(), "0021e010-97ab-46db-8600-f7604ab52f92")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(order)
}

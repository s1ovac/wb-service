package main

import (
	"context"
	"fmt"
	"log"

	"github.com/s1ovac/order-subscribe/internal/store/config"
	order "github.com/s1ovac/order-subscribe/internal/store/databases/order/db"
	"github.com/s1ovac/order-subscribe/internal/store/databases/postgresql"
)

func main() {
	//sb := subscribe.New()
	//err := sb.SubscribeToChannel()
	//if err != nil {
	//	log.Fatal(err)
	//}
	cfg := config.NewConfig()
	postgreSQL, err := postgresql.NewClient(context.TODO(), 3, cfg)
	if err != nil {
		log.Fatal(err)
	}
	rep := order.NewRepository(postgreSQL)
	order, err := rep.FindOne(context.TODO(), "b42062ec-689a-4561-8c90-d73929173652")
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(order)
}

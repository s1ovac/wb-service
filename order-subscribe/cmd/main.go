package main

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/s1ovac/order-subscribe/internal/store/config"
	"github.com/s1ovac/order-subscribe/internal/store/databases/order"
	"github.com/s1ovac/order-subscribe/internal/store/databases/postgresql"
	"github.com/s1ovac/order-subscribe/internal/subscribe"
)

func main() {
	sb := subscribe.New()
	err := sb.SubscribeToChannel()
	if err != nil {
		log.Fatal(err)
	}
	newOrder := order.Order{
		TrackNumber:       "WBILMTESTTRACK",
		Entry:             "WBIL",
		Locale:            "en",
		InternalSignature: "",
		CustomerID:        "test",
		DeliveryService:   "meest",
		ShardKey:          "9",
		SmID:              99,
		DateCreated:       time.Now(),
		OofShard:          "1",
	}
	cfg := config.NewConfig()
	postgreSQL, err := postgresql.NewClient(context.TODO(), 3, cfg)
	if err != nil {
		log.Fatal(err)
	}
	rep := order.NewRepository(postgreSQL)
	if err = rep.Create(context.TODO(), &newOrder); err != nil {
		log.Fatal(err)
	}
	order, err := rep.FindOne(context.TODO(), newOrder.OrderUID)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(order)
}

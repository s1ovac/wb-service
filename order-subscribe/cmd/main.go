package main

import (
	"log"

	"github.com/s1ovac/order-subscribe/internal/subscribe"
)

func main() {
	sb := subscribe.New()
	err := sb.SubscribeToChannel()
	if err != nil {
		log.Fatal(err)
	}
}

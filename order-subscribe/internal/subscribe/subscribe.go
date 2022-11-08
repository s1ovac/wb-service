package subscribe

import (
	"encoding/json"
	"fmt"

	"github.com/nats-io/stan.go"
	"github.com/s1ovac/order-subscribe/internal/store/databases/order"
)

type Subscribe struct {
	clusterID string
	clientID  string
	channel   string
}

func New() *Subscribe {
	return &Subscribe{
		clusterID: "test-cluster",
		clientID:  "order-suscriber",
		channel:   "order-notification",
	}
}

func (sb *Subscribe) SubscribeToChannel() error {
	sc, err := stan.Connect(sb.clusterID, sb.clientID)
	if err != nil {
		return fmt.Errorf("problem with connecting to channel: %s", err)
	}
	defer sc.Close()

	_, err = sc.Subscribe(sb.channel, handleOrder, stan.StartWithLastReceived())
	if err != nil {
		return fmt.Errorf("problem with reading channel: %s", err)
	}
	return nil
}

func handleOrder(orderMsg *stan.Msg) {
	newOrder := order.Order{}

	if err := json.Unmarshal(orderMsg.Data, &newOrder); err != nil {
		return
	}

}

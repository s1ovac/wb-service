package subscribe

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/nats-io/stan.go"
	"github.com/s1ovac/order-subscribe/internal/pkg/logging"
	"github.com/s1ovac/order-subscribe/internal/store/databases/order"
	"github.com/sirupsen/logrus"
)

type Subscribe struct {
	clusterID string
	clientID  string
	channel   string
	logger    *logrus.Logger
}

func New(logger *logrus.Logger) *Subscribe {
	return &Subscribe{
		clusterID: "test-cluster",
		clientID:  "order-suscriber",
		channel:   "order-notification",
		logger:    logging.Init(),
	}
}

func (sb *Subscribe) SubscribeToChannel() (*order.Order, error) {
	sb.logger.Infof("Connecting to the channel with\nclusterID: %s clientID: %s\n", sb.clusterID, sb.clientID)
	sc, err := stan.Connect(sb.clusterID, sb.clientID)
	var newOrder order.Order
	if err != nil {
		return nil, fmt.Errorf("problem with connecting to channel: %s", err)
	}
	defer sc.Close()

	sub, err := sc.Subscribe(sb.channel, func(orderMsg *stan.Msg) {
		if err := json.Unmarshal(orderMsg.Data, &newOrder); err != nil {
			log.Fatalf("error occured parsing json file: %s", err)
		}
	}, stan.StartWithLastReceived())
	if err != nil {
		return nil, fmt.Errorf("problem with reading channel: %s", err)
	}
	sub.Unsubscribe()
	return &newOrder, nil
}

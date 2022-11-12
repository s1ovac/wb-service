package subscribe

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nats-io/stan.go"
	"github.com/s1ovac/order-subscribe/internal/store/databases/order"
	"github.com/sirupsen/logrus"
)

type Subscriber struct {
	clusterID  string
	clientID   string
	channel    string
	logger     *logrus.Logger
	repository order.Repository
	conn       *pgxpool.Pool
}

func New(logger *logrus.Logger, repository order.Repository, conn *pgxpool.Pool) *Subscriber {
	return &Subscriber{
		clusterID:  "test-cluster",
		clientID:   "order-suscriber",
		channel:    "order-notification",
		logger:     logger,
		repository: repository,
		conn:       conn,
	}
}

func (sb *Subscriber) SubscribeToChannel(ctx context.Context) error {
	sb.logger.Infof("Connecting to the channel with\nclusterID: %s clientID: %s\n", sb.clusterID, sb.clientID)
	sc, err := stan.Connect(sb.clusterID, sb.clientID)
	var newOrder order.Order
	if err != nil {
		return fmt.Errorf("problem with connecting to channel: %s", err)
	}
	defer sc.Close()

	_, err = sc.Subscribe(sb.channel, func(orderMsg *stan.Msg) {
		if err := json.Unmarshal(orderMsg.Data, &newOrder); err != nil {
			sb.logger.Warning("Cannot unmarshal data from nats-streaming-server")
		}
	}, stan.StartWithLastReceived())
	if err != nil {
		return fmt.Errorf("problem with reading channel: %s", err)
	}
	if err := sb.repository.Create(ctx, &newOrder, sb.conn); err != nil {
		sb.logger.Error("can't subscribe to nats-streaming-channel")
		return err
	}
	return nil
}

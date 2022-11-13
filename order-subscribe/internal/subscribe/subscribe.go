package subscribe

import (
	"context"
	"encoding/json"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/nats-io/stan.go"
	"github.com/s1ovac/order-subscribe/internal/store/databases/order"
	"github.com/sirupsen/logrus"
)

type Subscriber struct {
	ClusterID  string
	ClientID   string
	Channel    string
	logger     *logrus.Logger
	repository order.Repository
	conn       *pgxpool.Pool
}

func New(logger *logrus.Logger, repository order.Repository, conn *pgxpool.Pool) *Subscriber {
	return &Subscriber{
		ClusterID:  "test-cluster",
		ClientID:   "order-suscriber",
		Channel:    "order-notification",
		logger:     logger,
		repository: repository,
		conn:       conn,
	}
}

func (sb *Subscriber) CreateOrder(m *stan.Msg) {
	order := order.Order{}
	if err := json.Unmarshal(m.Data, &order); err != nil {
		sb.logger.Error(err)
		return
	}
	if err := sb.repository.Create(context.TODO(), &order, sb.conn); err != nil {
		sb.logger.Error(err)
		return
	}

}

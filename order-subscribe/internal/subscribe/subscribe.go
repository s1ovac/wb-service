package subscribe

import (
	"context"
	"encoding/json"

	"github.com/go-playground/validator/v10"
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
	ctx        context.Context
}

func New(logger *logrus.Logger, repository order.Repository, conn *pgxpool.Pool, ctx context.Context) *Subscriber {
	return &Subscriber{
		ClusterID:  "test-cluster",
		ClientID:   "order-suscriber",
		Channel:    "order-notification",
		logger:     logger,
		repository: repository,
		conn:       conn,
		ctx:        ctx,
	}
}

func (sb *Subscriber) CreateOrder(m *stan.Msg) {
	order := order.Order{}
	err := json.Unmarshal(m.Data, &order)
	if err != nil {
		sb.logger.Error(err)
		return
	}
	validate := validator.New()
	err = validate.Struct(order)
	if err != nil {
		sb.logger.Error(err)
		return
	}
	if err := sb.repository.Create(sb.ctx, &order, sb.conn); err != nil {
		sb.logger.Error(err)
		return
	}

}

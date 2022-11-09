package order

import (
	"context"

	"github.com/jackc/pgx/v4"
)

type Repository interface {
	Create(ctx context.Context, order *Order, conn *pgx.Conn) error
	FindOne(ctx context.Context, id string) (Order, error)
	FindAll(ctx context.Context) (o []Order, err error)
}

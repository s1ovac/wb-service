package order

import (
	"context"

	"github.com/jackc/pgx/v4/pgxpool"
)

type Repository interface {
	Create(ctx context.Context, order *Order, conn *pgxpool.Pool) error
	FindOne(ctx context.Context, id string) (Order, error)
	FindAll(ctx context.Context) (o []Order, err error)
}

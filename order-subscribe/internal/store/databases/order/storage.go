package order

import "context"

type Repository interface {
	Create(ctx context.Context, order *Order) error
	FindOne(ctx context.Context, id string) (Order, error)
	FindAll(ctx context.Context) (o []Order, err error)
}

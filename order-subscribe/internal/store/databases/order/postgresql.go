package order

import (
	"context"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/s1ovac/order-subscribe/internal/store/databases/postgresql"
)

type repository struct {
	client postgresql.CLient
}

func NewRepository(client postgresql.CLient) Repository {
	return &repository{
		client: client,
	}
}

func (r *repository) Create(ctx context.Context, order *Order) error {
	q := `
	INSERT INTO "order" (
		"track_number", 
		"entry", 
		"locale", 
		"internal_signature", 
		"customer_id", 
		"delivery_service", 
		"shardkey", 
		"sm_id", 
		"date_created", 
		"oof_shard" ) 
	VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10) 
	RETURNING "order_uid"
	`
	if err := r.client.QueryRow(ctx, q, order.TrackNumber,
		order.Entry,
		order.Locale,
		order.InternalSignature,
		order.CustomerID,
		order.DeliveryService,
		order.ShardKey,
		order.SmID,
		order.DateCreated,
		order.OofShard,
	).Scan(&order.OrderUID); err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok {
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail %s, Where %s, Code %s, SQLState %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			fmt.Println(newErr)
			return newErr
		}
		return err
	}
	return nil
}

func (r *repository) FindOne(ctx context.Context, id string) (Order, error) {
	q := `
	SELECT "order_uid", "track_number", "entry", "locale", "internal_signature", "customer_id", "delivery_service", "shardkey", "sm_id", "date_created", "oof_shard"
	FROM "order"
	WHERE order_uid = $1
	`
	var ord Order
	if err := r.client.QueryRow(ctx, q, id).Scan(&ord.OrderUID, &ord.TrackNumber, &ord.Entry, &ord.Locale, &ord.InternalSignature, &ord.CustomerID, &ord.DeliveryService, &ord.ShardKey, &ord.SmID, &ord.DateCreated, &ord.OofShard); err != nil {
		return Order{}, err
	}
	return ord, nil
}

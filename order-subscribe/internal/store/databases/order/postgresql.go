package order

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/jackc/pgconn"
	"github.com/s1ovac/order-subscribe/internal/store/databases/item"
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
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			pgErr = err.(*pgconn.PgError)
			newErr := fmt.Errorf(fmt.Sprintf("SQL Error: %s, Detail %s, Where %s, Code %s, SQLState %s", pgErr.Message, pgErr.Detail, pgErr.Where, pgErr.Code, pgErr.SQLState()))
			log.Println(newErr)
			return newErr
		}
		return err
	}
	return nil
}

func (r *repository) FindOne(ctx context.Context, id string) (Order, error) {
	q := `
	SELECT 
		o."order_uid",
		o."track_number",
		o."entry",
		d."name",
		d."phone",
		d."zip",
		d."city",
		d."address",
		d."region",
		d."email",
		p."transaction",
		p."request_id",
		p."currency",
		p."provider",
		p."amount",
		p."payment_dt",
		p."bank",
		p."delivery_cost",
		p."goods_total",
		p."custom_fee",
		o."locale",
		o."internal_signature",
		o."customer_id",
		o."delivery_service", 
		o."shardkey", 
		o."sm_id", 
		o."date_created", 
		o."oof_shard"
	FROM 
		"order" AS o 
		JOIN "delivery" AS d ON o."order_uid" = d."order_id"
		JOIN "payment" AS p ON d."order_id" = p."order_id"
		RIGHT JOIN "item" AS i ON p."order_id" = i."order_id"
	WHERE o.order_uid = $1
	`
	var ord Order
	if err := r.client.QueryRow(ctx, q, id).Scan(
		&ord.OrderUID,
		&ord.TrackNumber,
		&ord.Entry,
		&ord.Delivery.Name,
		&ord.Delivery.Phone,
		&ord.Delivery.Zip,
		&ord.Delivery.City,
		&ord.Delivery.Address,
		&ord.Delivery.Region,
		&ord.Delivery.Email,
		&ord.Payment.Transaction,
		&ord.Payment.RequestID,
		&ord.Payment.Currency,
		&ord.Payment.Provider,
		&ord.Payment.Amount,
		&ord.Payment.PaymentDT,
		&ord.Payment.Bank,
		&ord.Payment.DeliveryCost,
		&ord.Payment.GoodsTotal,
		&ord.Payment.CustomFee,
		&ord.Locale,
		&ord.InternalSignature,
		&ord.CustomerID,
		&ord.DeliveryService,
		&ord.ShardKey,
		&ord.SmID,
		&ord.DateCreated,
		&ord.OofShard,
	); err != nil {
		return Order{}, err
	}
	iq := `
		SELECT 
			"id", 
			"chrt_id",
			"track_number",
			"price",
			"rid",
			"name",
			"sale",
			"size",
			"total_price",
			"nm_id",
			"brand",
			"status"
		FROM "item" AS i
		WHERE "order_id" = $1
		`
	itemRows, err := r.client.Query(ctx, iq, ord.OrderUID)
	if err != nil {
		return Order{}, err
	}
	items := make([]item.Item, 0)
	for itemRows.Next() {
		var item item.Item
		err = itemRows.Scan(
			&item.ID,
			&item.ChrtID,
			&item.TrackNumber,
			&item.Price,
			&item.Rid,
			&item.Name,
			&item.Sale,
			&item.Size,
			&item.TotalPrice,
			&item.NmID,
			&item.Brand,
			&item.Status,
		)
		if err != nil {
			return Order{}, err
		}
		items = append(items, item)
	}
	ord.Items = items
	return ord, nil
}

func (r *repository) FindAll(ctx context.Context) (o []Order, err error) {
	q := `
	SELECT 
		o."order_uid",
		o."track_number",
		o."entry",
		d."name",
		d."phone",
		d."zip",
		d."city",
		d."address",
		d."region",
		d."email",
		p."transaction",
		p."request_id",
		p."currency",
		p."provider",
		p."amount",
		p."payment_dt",
		p."bank",
		p."delivery_cost",
		p."goods_total",
		p."custom_fee",
		o."locale",
		o."internal_signature",
		o."customer_id",
		o."delivery_service", 
		o."shardkey", 
		o."sm_id", 
		o."date_created", 
		o."oof_shard"
	FROM 
		"order" AS o 
		JOIN "delivery" AS d ON o."order_uid" = d."order_id"
		JOIN "payment" AS p ON d."order_id" = p."order_id"
	`
	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	orders := make([]Order, 0)
	for rows.Next() {
		var ord Order
		err = rows.Scan(
			&ord.OrderUID,
			&ord.TrackNumber,
			&ord.Entry,
			&ord.Delivery.Name,
			&ord.Delivery.Phone,
			&ord.Delivery.Zip,
			&ord.Delivery.City,
			&ord.Delivery.Address,
			&ord.Delivery.Region,
			&ord.Delivery.Email,
			&ord.Payment.Transaction,
			&ord.Payment.RequestID,
			&ord.Payment.Currency,
			&ord.Payment.Provider,
			&ord.Payment.Amount,
			&ord.Payment.PaymentDT,
			&ord.Payment.Bank,
			&ord.Payment.DeliveryCost,
			&ord.Payment.GoodsTotal,
			&ord.Payment.CustomFee,
			&ord.Locale,
			&ord.InternalSignature,
			&ord.CustomerID,
			&ord.DeliveryService,
			&ord.ShardKey,
			&ord.SmID,
			&ord.DateCreated,
			&ord.OofShard,
		)
		if err != nil {
			return nil, err
		}

		iq := `
		SELECT 
			"id", 
			"chrt_id",
			"track_number",
			"price",
			"rid",
			"name",
			"sale",
			"size",
			"total_price",
			"nm_id",
			"brand",
			"status"
		FROM "item" AS i
		WHERE "order_id" = $1
		`
		itemRows, err := r.client.Query(ctx, iq, ord.OrderUID)
		if err != nil {
			return nil, err
		}
		items := make([]item.Item, 0)
		for itemRows.Next() {
			var item item.Item

			err = itemRows.Scan(
				&item.ID,
				&item.ChrtID,
				&item.TrackNumber,
				&item.Price,
				&item.Rid,
				&item.Name,
				&item.Sale,
				&item.Size,
				&item.TotalPrice,
				&item.NmID,
				&item.Brand,
				&item.Status,
			)
			if err != nil {
				return nil, err
			}
			items = append(items, item)
		}
		ord.Items = items

		orders = append(orders, ord)

	}
	return orders, nil
}

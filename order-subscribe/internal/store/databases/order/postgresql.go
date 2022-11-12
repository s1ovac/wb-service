package order

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/s1ovac/order-subscribe/internal/pkg/logging"
	"github.com/s1ovac/order-subscribe/internal/store/databases/item"
	"github.com/s1ovac/order-subscribe/internal/store/databases/postgresql"
	"github.com/sirupsen/logrus"
)

type repository struct {
	client postgresql.CLient
	logger *logrus.Logger
}

func NewRepository(client postgresql.CLient) Repository {
	return &repository{
		client: client,
		logger: logging.Init(),
	}
}

func (r *repository) Create(ctx context.Context, order *Order, conn *pgxpool.Pool) error {
	var (
		qOrder = `
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
		qDelivery = `
		INSERT INTO "delivery" (
			"order_id",
			"name", 
			"phone", 
			"zip", 
			"city", 
			"address", 
			"region", 
			"email"
			) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
		`
		qPayment = `
		INSERT INTO "payment" (
			"order_id",
			"transaction", 
			"request_id", 
			"currency", 
			"provider", 
			"amount", 
			"payment_dt", 
			"bank",
			"delivery_cost",
			"goods_total",
			"custom_fee"
			) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11)
		`
		qItem = `
		INSERT INTO "item" (
			"order_id",
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
			) 
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10, $11, $12) `
	)
	tx, err := conn.Begin(ctx)
	if err != nil {
		return err
	}
	r.logger.Info(fmt.Sprintf("SQL Query: %s", formatQuery(qOrder)))
	err = tx.QueryRow(
		ctx,
		qOrder,
		order.TrackNumber,
		order.Entry,
		order.Locale,
		order.InternalSignature,
		order.CustomerID,
		order.DeliveryService,
		order.ShardKey,
		order.SmID,
		order.DateCreated,
		order.OofShard,
	).Scan(&order.OrderUID)
	if err != nil {
		return err
	}
	r.logger.Info(fmt.Sprintf("SQL Query: %s", formatQuery(qDelivery)))
	_, err = tx.Exec(
		ctx,
		qDelivery,
		order.OrderUID,
		order.Delivery.Name,
		order.Delivery.Phone,
		order.Delivery.Zip,
		order.Delivery.City,
		order.Delivery.Address,
		order.Delivery.Region,
		order.Delivery.Email,
	)
	if err != nil {
		return err
	}
	r.logger.Info(fmt.Sprintf("SQL Query: %s", formatQuery(qPayment)))
	_, err = tx.Exec(
		ctx,
		qPayment,
		order.OrderUID,
		order.Payment.Transaction,
		order.Payment.RequestID,
		order.Payment.Currency,
		order.Payment.Provider,
		order.Payment.Amount,
		order.Payment.PaymentDT,
		order.Payment.Bank,
		order.Payment.DeliveryCost,
		order.Payment.GoodsTotal,
		order.Payment.CustomFee,
	)
	if err != nil {
		return err
	}
	r.logger.Info(fmt.Sprintf("SQL Query: %s", formatQuery(qItem)))
	for _, it := range order.Items {
		_, err = tx.Exec(
			ctx,
			qItem,
			order.OrderUID,
			it.ChrtID,
			order.TrackNumber,
			it.Price,
			it.Rid,
			it.Name,
			it.Sale,
			it.Size,
			it.TotalPrice,
			it.NmID,
			it.Brand,
			it.Status,
		)
		if err != nil {
			return err
		}
	}
	tx.Commit(ctx)

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
	r.logger.Info(fmt.Sprintf("SQL Query: %s", formatQuery(q)))
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
	r.logger.Info(fmt.Sprintf("SQL Query: %s", formatQuery(iq)))
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
		d."order_id",
		d."name",
		d."phone",
		d."zip",
		d."city",
		d."address",
		d."region",
		d."email",
		p."order_id",
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
	r.logger.Info(fmt.Sprintf("SQL Query: %s", formatQuery(q)))
	rows, err := r.client.Query(ctx, q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	orders := make([]Order, 0)
	for rows.Next() {
		var ord Order
		err = rows.Scan(
			&ord.OrderUID,
			&ord.TrackNumber,
			&ord.Entry,
			&ord.Delivery.OrderID,
			&ord.Delivery.Name,
			&ord.Delivery.Phone,
			&ord.Delivery.Zip,
			&ord.Delivery.City,
			&ord.Delivery.Address,
			&ord.Delivery.Region,
			&ord.Delivery.Email,
			&ord.Payment.OrderID,
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
			"order_id", 
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
		r.logger.Info(fmt.Sprintf("SQL Query: %s", formatQuery(iq)))
		itemRows, err := r.client.Query(ctx, iq, ord.OrderUID)
		if err != nil {
			return nil, err
		}
		defer itemRows.Close()
		items := make([]item.Item, 0)
		for itemRows.Next() {
			var item item.Item

			err = itemRows.Scan(
				&item.ID,
				&item.OrderID,
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

func formatQuery(q string) string {
	return strings.ReplaceAll(strings.ReplaceAll(q, "\t", ""), "\n", " ")
}

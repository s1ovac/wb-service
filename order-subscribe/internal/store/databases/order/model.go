package order

import (
	"time"

	"github.com/s1ovac/order-subscribe/internal/store/databases/delivery"
	"github.com/s1ovac/order-subscribe/internal/store/databases/item"
	"github.com/s1ovac/order-subscribe/internal/store/databases/payment"
)

type Order struct {
	OrderUID          string            `json:"order_uid"`
	TrackNumber       string            `json:"track_number" validate:"required"`
	Entry             string            `json:"entry" validate:"required"`
	Delivery          delivery.Delivery `json:"delivery" validate:"required"`
	Payment           payment.Payment   `json:"payment" validate:"required"`
	Items             []item.Item       `json:"items" validate:"required"`
	Locale            string            `json:"locale" validate:"required"`
	InternalSignature string            `json:"internal_signature" validate:"required"`
	CustomerID        string            `json:"customer_id" validate:"required"`
	DeliveryService   string            `json:"delivery_service" validate:"required"`
	ShardKey          string            `json:"shardkey" validate:"required"`
	SmID              int               `json:"sm_id" validate:"required"`
	DateCreated       time.Time         `json:"date_created" validate:"required"`
	OofShard          string            `json:"oof_shard" validate:"required"`
}

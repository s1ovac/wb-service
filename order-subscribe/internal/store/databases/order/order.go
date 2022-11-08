package order

import (
	"time"

	"github.com/s1ovac/order-subscribe/internal/store/databases/delivery"
	"github.com/s1ovac/order-subscribe/internal/store/databases/item"
	"github.com/s1ovac/order-subscribe/internal/store/databases/payment"
)

type Order struct {
	OrderUID          string            `json:"order_uid"`
	TrackNumber       string            `json:"track_number"`
	Entry             string            `json:"entry"`
	Delivery          delivery.Delivery `json:"delivery"`
	Payment           payment.Payment   `json:"payment"`
	Items             []item.Item       `json:"items"`
	Locale            string            `json:"locale"`
	InternalSignature string            `json:"internal_signature"`
	CustomerID        string            `json:"customer_id"`
	DeliveryService   string            `json:"delivery_service"`
	ShardKey          string            `json:"shardkey"`
	SmID              int               `json:"sm_id"`
	DateCreated       time.Time         `json:"date_created"`
	OofShard          string            `json:"oof_shard"`
}

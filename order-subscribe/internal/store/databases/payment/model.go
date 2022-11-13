package payment

type Payment struct {
	OrderID      string  `json:"order_id"`
	Transaction  string  `json:"transaction" validate:"required"`
	RequestID    string  `json:"request_id" validate:"required"`
	Currency     string  `json:"currency" validate:"required"`
	Provider     string  `json:"provider" validate:"required"`
	Amount       int     `json:"amount" validate:"required"`
	PaymentDT    int     `json:"payment_dt" validate:"required"`
	Bank         string  `json:"bank" validate:"required"`
	DeliveryCost float64 `json:"delivery_cost" validate:"required"`
	GoodsTotal   int     `json:"goods_total" validate:"required"`
	CustomFee    float64 `json:"custom_fee"`
}

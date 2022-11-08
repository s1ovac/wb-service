package model

type Payment struct {
	OrderID      string  `json:""`
	Transaction  string  `json:""`
	RequestID    string  `json:""`
	Currency     string  `json:""`
	Provider     string  `json:""`
	Amount       string  `json:""`
	PaymentDT    int     `json:""`
	Bank         string  `json:""`
	DeliveryCost float64 `json:""`
	GoodsTotal   int     `json:""`
	CustomFee    float64 `json:""`
}

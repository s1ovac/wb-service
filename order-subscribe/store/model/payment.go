package model

type Payment struct {
	OrderID      string
	Transaction  string
	RequestID    string
	Currency     string
	Provider     string
	Amount       string
	PaymentDT    int
	Bank         string
	DeliveryCost float64
	GoodsTotal   int
	CustomFee    float64
}

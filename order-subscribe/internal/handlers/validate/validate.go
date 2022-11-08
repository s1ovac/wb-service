package validate

import (
	"encoding/json"
	"fmt"

	"github.com/nats-io/stan.go"
	"github.com/s1ovac/order-subscribe/internal/store/databases/order"
)

type Data struct {
	model order.Order
}

func NewData() *Data {
	return &Data{}
}

func (d *Data) Validate(orderMsg *stan.Msg) error {
	if err := json.Unmarshal(orderMsg.Data, &d.model); err != nil {
		return fmt.Errorf("error occured while parsing json data: %s", err)
	}
	return nil
}

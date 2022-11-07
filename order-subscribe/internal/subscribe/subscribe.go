package subscribe

import (
	"fmt"
	"time"

	"github.com/nats-io/stan.go"
)

type Subscribe struct {
	clusterID string
	clientID  string
	channel   string
}

func New() *Subscribe {
	return &Subscribe{
		clusterID: "test-cluster",
		clientID:  "order-suscriber",
		channel:   "order-notification",
	}
}

func (sb *Subscribe) SubscribeToChannel() error {
	sc, err := stan.Connect(sb.clusterID, sb.clientID)
	if err != nil {
		return fmt.Errorf("problem with connecting to channel: %s", err)
	}
	defer sc.Close()

	sub, err := sc.Subscribe(sb.channel, func(m *stan.Msg) {
		fmt.Printf("Received a message: %s\n", string(m.Data))
	})

	if err != nil {
		return fmt.Errorf("problem with reading channel: %s", err)
	}
	time.Sleep(30 * time.Second)
	sub.Unsubscribe()
	return nil
}

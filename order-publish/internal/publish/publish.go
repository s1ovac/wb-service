package publish

import (
	"bufio"
	"fmt"
	"os"
	"strings"

	"github.com/nats-io/stan.go"
)

type Publish struct {
	clusterID string
	clientID  string
	channel   string
}

func New() *Publish {
	return &Publish{
		clusterID: "test-cluster",
		clientID:  "order-publisher",
		channel:   "order-notification",
	}
}

func (pb *Publish) DropMessage() error {
	sc, err := stan.Connect(pb.clusterID, pb.clientID)
	if err != nil {
		return err
	}
	defer sc.Close()
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Println("Enter path to JSON file")
		path, err := reader.ReadString('\n')
		if err != nil {
			return fmt.Errorf("incorrect path to JSON file: %s", err)
		}
		path = strings.TrimSpace(path)
		if path != "" {
			file, err := os.ReadFile(path)
			if err != nil {
				return fmt.Errorf("problem with reading file: %s", err)
			}
			if err = sc.Publish(pb.channel, file); err != nil {
				return fmt.Errorf("problem with publish: %s", err)
			}
		}
		return nil
	}
}

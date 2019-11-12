package mqhub

import (
	"fmt"
	"testing"
)

func TestKafkaSubscribe(t *testing.T) {
	k := &KafkaSubscriber{}

	k.broker_list = []string{"114.212.87.225:32590"}
	k.group = "testworkflowEngine127.0.0.1"

	valueChan, _ := k.SubscribeTo("event-source")

	for {
		select {
		case value := <-valueChan:
			fmt.Printf("%% Message :\n%s\n", value)
		}
	}
}

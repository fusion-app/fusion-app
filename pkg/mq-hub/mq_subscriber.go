package mqhub

import (
	"fmt"
	"github.com/Shopify/sarama"
	cluster "github.com/bsm/sarama-cluster"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
)

type KafkaSubscriber struct {
	broker_list []string
	group       string
	// topics		[]string

	// valueChan	chan[1000] string
}

// kafka config init
// return error
func (k *KafkaSubscriber) Init(config interface{}) error {
	// k.broker_list = config.broker_list
	// k.topics = config.topics
	// k.group = config.group
	return nil
}

func (k *KafkaSubscriber) SubscribeTo(topic string) (<-chan interface{}, error) {
	config := cluster.NewConfig()
	config.Consumer.Return.Errors = true
	config.Group.Return.Notifications = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	var topics []string
	topics = append(topics, topic)

	consumer, err := cluster.NewConsumer(k.broker_list, k.group, topics, config)

	valueChan := make(chan interface{}, 1000)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create consumer: %s\n", err)
		panic(err)
	}
	// defer consumer.Close()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, os.Interrupt)

	fmt.Printf("Created Consumer %v\n", consumer)

	go func() {
		for err := range consumer.Errors() {
			log.Printf("Error: %s\n", err.Error())
		}
	}()

	go func() {
		for ntf := range consumer.Notifications() {
			log.Printf("Rebalanced: %+v\n", ntf)
		}
	}()

	go func() {
		for {
			select {
			case msg, ok := <-consumer.Messages():
				if ok {
					valueChan <- string(msg.Value)
					log.Printf("%s/%d/%d\t%s\t%s\n", msg.Topic, msg.Partition, msg.Offset, msg.Key, msg.Value)
					consumer.MarkOffset(msg, "")
				}
			case <-signals:
				return
			}
		}
	}()

	return valueChan, err
}

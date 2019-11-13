package subscriber

import (
	"encoding/json"
	jsonpatch "github.com/evanphx/json-patch"
	"github.com/fusion-app/fusion-app/pkg/client"
	"github.com/fusion-app/fusion-app/pkg/mq-hub"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
	"strings"
)

func NewSubscribeCommand() *cobra.Command {
	var broker string
	var group string
	var topic string
	var command = &cobra.Command{
		Use:      "subscribe FLAG",
		Short:    "subscribe a topic",
		Long:     `subscribe a topic`,
		Args:     cobra.MaximumNArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			if err := subscribeTopic(topic, broker, group); err != nil {
				log.Fatalf("Failed to deploy function: %v", err)
				os.Exit(1)
			}
		},
	}
	command.Flags().StringVar(&topic, "topic", "", "Specify topic")
	command.Flags().StringVar(&broker, "broker", "", "Specify broker")
	command.Flags().StringVar(&group, "group", "", "Specify group")
	return command
}

func subscribeTopic(topic string, broker string, group string) error {
	kafkaSubscriber := &mqhub.KafkaSubscriber{}
	_ = kafkaSubscriber.Init(strings.Split(broker, ","), group)
	valueChan, _ := kafkaSubscriber.SubscribeTo(topic)
	resourceClient, err := client.NewResources("fusion-app-resources")
	if err != nil {
		return err
	}
	for {
		select {
		case value := <-valueChan:
			msg := mqhub.Message{}
			bytes := []byte(value.(string))
			if err := json.Unmarshal(bytes, &msg); err != nil {
				return err
			}
			resource, err := resourceClient.Get(msg.Target.Name, v1.GetOptions{})
			if err != nil {
				return err
			}
			original, err := json.Marshal(resource.Labels)
			patchJson, err := json.Marshal(msg.UpdatePatch)
			patch, err := jsonpatch.DecodePatch(patchJson)
			if err != nil {
				return err
			}
			modified, err := patch.Apply(original)
			if err != nil {
				return err
			}
			newLabels := map[string]string{}
			err = json.Unmarshal(modified, &newLabels)
			resource.SetLabels(newLabels)
			resource, err = resourceClient.Update(resource)
			if err != nil {
				return err
			}
		}
	}
}
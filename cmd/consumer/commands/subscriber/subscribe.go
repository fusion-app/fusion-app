package subscriber

import (
	"context"
	"encoding/json"
	"fmt"
	jsonpatch "github.com/evanphx/json-patch"
	"github.com/fusion-app/fusion-app/pkg/apis"
	"github.com/fusion-app/fusion-app/pkg/apis/fusionapp/v1alpha1"
	"github.com/fusion-app/fusion-app/pkg/mq-hub"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/rest"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
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
	kubeConfig, err := config.GetConfig()
	if err != nil {
		return err
	}
	clientset, err := setupClient(kubeConfig)
	if err != nil {
		return fmt.Errorf("failed to setup kubernetes client: %v", err)
	}
	for {
		select {
		case value := <-valueChan:
			msg := mqhub.Message{}
			bytes := []byte(value.(string))
			if err := json.Unmarshal(bytes, &msg); err != nil {
				return err
			}
			resource := new(v1alpha1.Resource)
			err := clientset.Get(context.TODO(), client.ObjectKey{Namespace: msg.Target.Namespace,
				Name: msg.Target.Name}, resource)
			if errors.IsNotFound(err) {
				continue
			} else if err != nil {
				return err
			}
			original, err := json.Marshal(resource.Spec.Labels)
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
			in, out := &newLabels, &resource.Spec.Labels
			*out = make(map[string]string, len(*in))
			for key, val := range *in {
				(*out)[key] = val
			}
			err = clientset.Update(context.TODO(), resource)
			if err != nil {
				return err
			}
		}
	}
}

func setupClient(config *rest.Config) (client.Client, error) {
	scheme := runtime.NewScheme()
	for _, addToSchemeFunc := range []func(s *runtime.Scheme) error{
		apis.AddToScheme,
		corev1.AddToScheme,
	} {
		if err := addToSchemeFunc(scheme); err != nil {
			return nil, err
		}
	}

	clientset, err := client.New(config, client.Options{Scheme: scheme})
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

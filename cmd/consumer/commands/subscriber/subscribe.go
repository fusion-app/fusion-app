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
			log.Printf("Starting to process msg: %s", value.(string))
			if err := json.Unmarshal(bytes, &msg); err != nil {
				log.Printf("Failed to parse message: %v", err)
				continue
			}
			var obj runtime.Object
			var original []byte
			if msg.Target.Kind ==  "Resource" {
				resource := new(v1alpha1.Resource)
				err := clientset.Get(context.TODO(), client.ObjectKey{Namespace: msg.Target.Namespace,
					Name: msg.Target.Name}, resource)
				if errors.IsNotFound(err) {
					continue
				} else if err != nil {
					log.Printf("Failed to get resource: %v", err)
					continue
				}
				if resource.Spec.Labels == nil {
					resource.Spec.Labels = map[string]string{}
				}
				original, err = json.Marshal(resource.Spec.Labels)
				obj = resource
			} else if msg.Target.Kind ==  "FusionApp" {
				app := new(v1alpha1.FusionApp)
				err := clientset.Get(context.TODO(), client.ObjectKey{Namespace: msg.Target.Namespace,
					Name: msg.Target.Name}, app)
				if errors.IsNotFound(err) {
					continue
				} else if err != nil {
					log.Printf("Failed to get fusionApp: %v", err)
					continue
				}
				if app.Spec.Labels == nil {
					app.Spec.Labels = map[string]string{}
				}
				original, err = json.Marshal(app.Spec.Labels)
				obj = app
			} else if msg.Target.Kind ==  "FusionAppInstance" {
				appInstance := new(v1alpha1.FusionAppInstance)
				err := clientset.Get(context.TODO(), client.ObjectKey{Namespace: msg.Target.Namespace,
					Name: msg.Target.Name}, appInstance)
				if errors.IsNotFound(err) {
					continue
				} else if err != nil {
					log.Printf("Failed to get fusionAppInstance: %v", err)
					continue
				}
				if appInstance.Spec.Labels == nil {
					appInstance.Spec.Labels = map[string]string{}
				}
				original, err = json.Marshal(appInstance.Spec.Labels)
				obj = appInstance
			} else {
				log.Printf("No such kind: %s", msg.Target.Kind)
				continue
			}
			patchJson, err := json.Marshal(msg.UpdatePatch)
			patch, err := jsonpatch.DecodePatch(patchJson)
			if err != nil {
				log.Printf("Failed to parse patch: %v", err)
				continue
			}
			log.Printf("patchJson: %s", string(patchJson))
			log.Printf("original: %s", string(original))
			modified, err := patch.Apply(original)
			if err != nil {
				log.Printf("Failed to apply patch: %v", err)
				continue
			}
			log.Printf("modified: %s", string(modified))
			newLabels := map[string]string{}
			err = json.Unmarshal(modified, &newLabels)
			if msg.Target.Kind ==  "Resource" {
				resource := obj.(*v1alpha1.Resource)
				in, out := &newLabels, &resource.Spec.Labels
				*out = make(map[string]string, len(*in))
				for key, val := range *in {
					(*out)[key] = val
				}
				err = clientset.Update(context.TODO(), resource)
				if err != nil {
					log.Printf("Failed to update resource: %v", err)
					continue
				}
			} else if msg.Target.Kind ==  "FusionApp" {
				app := obj.(*v1alpha1.FusionApp)
				in, out := &newLabels, &app.Spec.Labels
				*out = make(map[string]string, len(*in))
				for key, val := range *in {
					(*out)[key] = val
				}
				err = clientset.Update(context.TODO(), app)
				if err != nil {
					log.Printf("Failed to update fusionApp: %v", err)
					continue
				}
			} else if msg.Target.Kind ==  "FusionAppInstance" {
				appInstance := obj.(*v1alpha1.FusionAppInstance)
				in, out := &newLabels, &appInstance.Spec.Labels
				*out = make(map[string]string, len(*in))
				for key, val := range *in {
					(*out)[key] = val
				}
				err = clientset.Update(context.TODO(), appInstance)
				if err != nil {
					log.Printf("Failed to update fusionAppInstance: %v", err)
					continue
				}
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

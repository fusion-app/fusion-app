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
			var originalLabels, originalStatus []byte
			var modifiedLables, modifiedStatus []byte
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
				originalLabels, err = json.Marshal(resource.Spec.Labels)
				originalStatus, err = json.Marshal(resource.Status)
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
				originalLabels, err = json.Marshal(app.Spec.Labels)
				originalStatus, err = json.Marshal(app.Status)
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
				originalLabels, err = json.Marshal(appInstance.Spec.Labels)
				originalStatus, err = json.Marshal(appInstance.Status)
				obj = appInstance
			} else {
				log.Printf("No such kind: %s", msg.Target.Kind)
				continue
			}
			if msg.LabelsPatch != nil && len(msg.LabelsPatch) > 0 {
				labelsPatchJson, err := json.Marshal(msg.LabelsPatch)
				labelsPatch, err := jsonpatch.DecodePatch(labelsPatchJson)
				if err != nil {
					log.Printf("Failed to parse labelsPatch: %v", err)
					continue
				}
				log.Printf("labelsPatchJson: %s", string(labelsPatchJson))
				log.Printf("originalLabels: %s", string(originalLabels))
				modifiedLables, err = labelsPatch.Apply(originalLabels)
				if err != nil {
					log.Printf("Failed to apply labelsPatch: %v", err)
					continue
				}
				log.Printf("modifiedLabels: %s", string(modifiedLables))
			}
			if msg.StatusPatch != nil && len(msg.StatusPatch) > 0 {
				statusPatchJson, err := json.Marshal(msg.StatusPatch)
				statusPatch, err := jsonpatch.DecodePatch(statusPatchJson)
				if err != nil {
					log.Printf("Failed to parse statusPatch: %v", err)
					continue
				}
				log.Printf("statusPatchJson: %s", string(statusPatchJson))
				log.Printf("originalStatus: %s", string(originalStatus))
				modifiedStatus, err = statusPatch.Apply(originalStatus)
				if err != nil {
					log.Printf("Failed to apply statusPatch: %v", err)
					continue
				}
				log.Printf("modifiedStatus: %s", string(modifiedStatus))
			}
			if msg.Target.Kind ==  "Resource" {
				resource := obj.(*v1alpha1.Resource)
				if modifiedLables != nil {
					newLabels := map[string]string{}
					err := json.Unmarshal(modifiedLables, &newLabels)
					if err != nil {
						log.Printf("Failed to modify Labels: %v", err)
						continue
					}
					in, out := &newLabels, &resource.Spec.Labels
					*out = make(map[string]string, len(*in))
					for key, val := range *in {
						(*out)[key] = val
					}
				}
				if modifiedStatus != nil {
					newStatus := v1alpha1.ResourceStatus{}
					err := json.Unmarshal(modifiedStatus, &newStatus)
					if err != nil {
						log.Printf("Failed to modify Status: %v", err)
						continue
					}
					in, out := &newStatus, &resource.Status
					in.DeepCopyInto(out)
				}
				err = clientset.Update(context.TODO(), resource)
				if err != nil {
					log.Printf("Failed to update resource: %v", err)
					continue
				}
			} else if msg.Target.Kind ==  "FusionApp" {
				app := obj.(*v1alpha1.FusionApp)
				if modifiedLables != nil {
					newLabels := map[string]string{}
					err := json.Unmarshal(modifiedLables, &newLabels)
					if err != nil {
						log.Printf("Failed to modify Labels: %v", err)
						continue
					}
					in, out := &newLabels, &app.Spec.Labels
					*out = make(map[string]string, len(*in))
					for key, val := range *in {
						(*out)[key] = val
					}
				}
				if modifiedStatus != nil {
					newStatus := v1alpha1.FusionAppStatus{}
					err := json.Unmarshal(modifiedStatus, &newStatus)
					if err != nil {
						log.Printf("Failed to modify Status: %v", err)
						continue
					}
					in, out := &newStatus, &app.Status
					in.DeepCopyInto(out)
				}
				err = clientset.Update(context.TODO(), app)
				if err != nil {
					log.Printf("Failed to update fusionApp: %v", err)
					continue
				}
			} else if msg.Target.Kind ==  "FusionAppInstance" {
				appInstance := obj.(*v1alpha1.FusionAppInstance)
				if modifiedLables != nil {
					newLabels := map[string]string{}
					err := json.Unmarshal(modifiedLables, &newLabels)
					if err != nil {
						log.Printf("Failed to modify Labels: %v", err)
						continue
					}
					in, out := &newLabels, &appInstance.Spec.Labels
					*out = make(map[string]string, len(*in))
					for key, val := range *in {
						(*out)[key] = val
					}
				}
				if modifiedStatus != nil {
					newStatus := v1alpha1.FusionAppInstanceStatus{}
					err := json.Unmarshal(modifiedStatus, &newStatus)
					if err != nil {
						log.Printf("Failed to modify Status: %v", err)
						continue
					}
					in, out := &newStatus, &appInstance.Status
					in.DeepCopyInto(out)
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

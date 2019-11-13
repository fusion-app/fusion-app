package resource

import (
	"crypto/sha1"
	"fmt"
	fusionappv1alpha1 "github.com/fusion-app/fusion-app/pkg/apis/fusionapp/v1alpha1"
	"io"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
)

const (
	probeImage             = "registry.njuics.cn/fusion-app/http-prober:201910122153"
	probeCommand           = "/usr/local/bin/ResourceProbeExample"

	consumerImage          = "registry.njuics.cn/fusion-app/consumer:2019111310"
	consumerCommand        = "/usr/local/bin/fusionapp-consumer"
	consumerServiceAccount = "fusion-app"

	defaultMqAddress       = "114.212.87.225:32015"
	defaultBroker          = "114.212.87.225:32590"

	EnvMqAdress            = "MQ_ADRESS"
	EnvBroker              = "BROKER"
)

// newPodForProbe returns a probe pod with the same name/namespace as the cr
func newPodForProbe(resource *fusionappv1alpha1.Resource) *corev1.Pod {
	mqAddress := os.Getenv(EnvMqAdress)
	if len(mqAddress) == 0 {
		mqAddress = defaultMqAddress
	}
	topic := topicForResource(resource)
	commands := []string{probeCommand, "--mq-address", mqAddress, "--mq-topic", topic, "--crd-namespace",
		resource.Namespace}
	commands = append(commands, resource.Spec.ProbeArgs...)
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resource.Name + "-probe-pod",
			Namespace: resource.Namespace,
			Labels:    DefaultLabels(resource),
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    resource.Name,
					Image:   probeImage,
					Command: commands,
				},
			},
		},
	}
}

// newPodForConsumer returns a kafka consumer pod with the same name/namespace as the cr
func newPodForConsumer(resource *fusionappv1alpha1.Resource) *corev1.Pod {
	broker := os.Getenv(EnvBroker)
	if len(broker) == 0 {
		broker = defaultBroker
	}
	topic := topicForResource(resource)
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resource.Name + "-consumer-pod",
			Namespace: resource.Namespace,
			Labels:    DefaultLabels(resource),
		},
		Spec: corev1.PodSpec{
			ServiceAccountName: consumerServiceAccount,
			Containers: []corev1.Container{
				{
					Name:    resource.Name,
					Image:   consumerImage,
					Command: []string{consumerCommand, "subscribe", "--topic",
						topic, "--broker", broker, "--group", topic},
					ImagePullPolicy: corev1.PullAlways,
				},
			},
		},
	}
}

func topicForResource(resource *fusionappv1alpha1.Resource) string {
	resourceAndKind := fmt.Sprintf("%s-%s-", resource.Name, string(resource.Spec.ResourceKind))
	sha1Hash := sha1.New()
	_, _ = io.WriteString(sha1Hash, resourceAndKind)
	topic := fmt.Sprintf("%s%s", resourceAndKind, sha1Hash.Sum(nil))
	return topic
}
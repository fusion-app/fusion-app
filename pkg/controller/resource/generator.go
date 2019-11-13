package resource

import (
	fusionappv1alpha1 "github.com/fusion-app/fusion-app/pkg/apis/fusionapp/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
)

const (
	probeImage             = "registry.njuics.cn/fusion-app/http-prober:201910122153"
	probeCommand           = "/usr/local/bin/ResourceProbeExample"

	topic                  = "resource-event-source"

	defaultMqAddress       = "114.212.87.225:32015"
	defaultBroker          = "114.212.87.225:32590"

	EnvMqAdress            = "MQ_ADRESS"
)

// newPodForProbe returns a probe pod with the same name/namespace as the cr
func newPodForProbe(resource *fusionappv1alpha1.Resource) *corev1.Pod {
	mqAddress := os.Getenv(EnvMqAdress)
	if len(mqAddress) == 0 {
		mqAddress = defaultMqAddress
	}
	commands := []string{probeCommand, "--mq-address", mqAddress, "--mq-topic", topic, "--crd-namespace",
		resource.Namespace, "--crd-name", resource.Name, "--crd-kind", string(resource.Spec.ResourceKind)}
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

//// newPodForConsumer returns a kafka consumer pod with the same name/namespace as the cr
//func newPodForConsumer(resource *fusionappv1alpha1.Resource) *corev1.Pod {
//	broker := os.Getenv(EnvBroker)
//	if len(broker) == 0 {
//		broker = defaultBroker
//	}
//	topic := topicForResource(resource)
//	return &corev1.Pod{
//		ObjectMeta: metav1.ObjectMeta{
//			Name:      resource.Name + "-consumer-pod",
//			Namespace: resource.Namespace,
//			Labels:    DefaultLabels(resource),
//		},
//		Spec: corev1.PodSpec{
//			ServiceAccountName: consumerServiceAccount,
//			Containers: []corev1.Container{
//				{
//					Name:    resource.Name,
//					Image:   consumerImage,
//					Command: []string{consumerCommand, "subscribe", "--topic",
//						topic, "--broker", broker, "--group", topic},
//					ImagePullPolicy: corev1.PullAlways,
//				},
//			},
//		},
//	}
//}
//
//func topicForResource(resource *fusionappv1alpha1.Resource) string {
//	topic := fmt.Sprintf("%s-%s", resource.Name, resource.Namespace)
//	return topic
//}
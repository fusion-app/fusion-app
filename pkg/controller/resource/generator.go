package resource

import (
	fusionappv1alpha1 "github.com/fusion-app/fusion-app/pkg/apis/fusionapp/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8sresource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
)

const (
	defaultProbeImage      = "registry.cn-shanghai.aliyuncs.com/fusion-app/http-prober:resource-prober.201911211325"

	topic                  = "resource-event-source"

	defaultMqAddress       = "221.228.66.83:30595"    // "114.212.87.225:32015"

	EnvMqAdress            = "MQ_ADRESS"
)

// newDeployForProbe returns a probe deployment with the same name/namespace as the cr
func newDeployForProbe(resource *fusionappv1alpha1.Resource) *appsv1.Deployment {
	probeImage := defaultProbeImage
	if len(resource.Spec.ProbeImage) > 0 {
		probeImage = resource.Spec.ProbeImage
	}
	mqAddress := os.Getenv(EnvMqAdress)
	if len(mqAddress) == 0 {
		mqAddress = defaultMqAddress
	}
	args := []string{"--mq-address", mqAddress, "--mq-topic", topic, "--crd-namespace",
		resource.Namespace, "--crd-name", resource.Name, "--crd-kind", resource.Kind,
		"--crd-uid", string(resource.UID)}
	args = append(args, resource.Spec.ProbeArgs...)
	return &appsv1.Deployment {
		ObjectMeta: metav1.ObjectMeta{
			Name:      resource.Name + "-probe-deploy",
			Namespace: resource.Namespace,
			Labels:    DefaultLabels(resource),
		},
		Spec: appsv1.DeploymentSpec{
			Selector: metav1.SetAsLabelSelector(DefaultLabels(resource)),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:    DefaultLabels(resource),
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    resource.Name,
							Image:   probeImage,
							Args:    args,
							Resources: corev1.ResourceRequirements{
								Limits: corev1.ResourceList{
									corev1.ResourceCPU: k8sresource.MustParse("500m"),
									corev1.ResourceMemory: k8sresource.MustParse("2Gi"),
								},
							},
						},
					},
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
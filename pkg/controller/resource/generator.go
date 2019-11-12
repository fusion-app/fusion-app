package resource

import (
	fusionappv1alpha1 "github.com/fusion-app/fusion-app/pkg/apis/fusionapp/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

const (
	probeImage = "registry.cn-hangzhou.aliyuncs.com/njuicscn/"
	probeCommand = "ResourceProbeExample"

	consumerImage = "registry.cn-hangzhou.aliyuncs.com/njuicscn/"
	consumerCommand = "fusionapp-consumer"
)

// newPodForProbe returns a probe pod with the same name/namespace as the cr
func newPodForProbe(resource *fusionappv1alpha1.Resource) *corev1.Pod {
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
					Command: append([]string{probeCommand}, resource.Spec.ProbeArgs...),
				},
			},
		},
	}
}

// newPodForConsumer returns a kafka consumer pod with the same name/namespace as the cr
func newPodForConsumer(resource *fusionappv1alpha1.Resource) *corev1.Pod {
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      resource.Name + "-consumer-pod",
			Namespace: resource.Namespace,
			Labels:    DefaultLabels(resource),
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:    resource.Name,
					Image:   consumerImage,
					Command: []string{"sleep", "3600"},
				},
			},
		},
	}
}

package util

import (
	"context"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/labels"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func ResourcesForLabels(namespace string, labels labels.Set, c client.Client) ([]v1.Pod, error) {
	pods := &v1.PodList{}
	err := c.List(context.TODO(),
		client.InNamespace(namespace).
			MatchingLabels(labels), pods)

	if err != nil {
		return nil, err
	}
	return pods.Items, nil
}

func MappingResourcesBy(pods []v1.Pod) map[v1.PodPhase]int {
	result := make(map[v1.PodPhase]int)
	for _, pod := range pods {
		if len(pod.Status.Phase) == 0 {
			continue
		}
		if _, ok := result[pod.Status.Phase]; !ok {
			result[pod.Status.Phase] = 1
		} else {
			result[pod.Status.Phase]++
		}
	}
	return result
}

package resource

import (
	"encoding/json"
	"fmt"
	fusionappv1alpha1 "github.com/fusion-app/fusion-app/pkg/apis/fusionapp/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8sresource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const (
	defaultProbeImage      = "registry.cn-shanghai.aliyuncs.com/fusion-app/http-prober:resource-prober.202010271442"

	defaultMSImage         = "registry.cn-hangzhou.aliyuncs.com/tangcong/airpurifier_service:v2"

	defaultMQBrokerHost    = "message-broker.fusion-app.svc.cluster.local"

	defaultSpringBootAdminHost = "spring-boot-admin-svc.fusion-app.svc.cluster.local:8000"

	EnvSpringBootAdminHost = "SPRING_BOOT_ADMIN_HOST"

	EnvMQBrokerHost        = "MQ_BROKER_HOST"

	envProbeImage          = "RESOURCE_PROBE_IMAGE"

	defaultMSPort    int32 = 8081

	defaultMSPortName      = "http"

	defaultPatcherPath = "/root/fusion-app"

	defaultPathcerFileName = "patcherconfig.json"
)

func newConfigmapForPatcher(resource *fusionappv1alpha1.Resource) *corev1.ConfigMap {
	jsonData, _ := json.Marshal(fusionappv1alpha1.PatcherConfig{Patchers: resource.Spec.ProbeSpec.Patchers})
	configmap := &corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name: resource.Name + "-cm",
			Namespace: resource.Namespace,
		},
		Data: map[string]string{
			"patcherconfig" : string(jsonData),
		},
	}
	return configmap
}

// newDeployForProbeAndMS returns a probe deployment with the same name/namespace as the cr
func newDeployForProbeAndMS(resource *fusionappv1alpha1.Resource) *appsv1.Deployment {
	MSImage := defaultMSImage
	probeImage := defaultProbeImage
	port := defaultMSPort
	if resource.Spec.ConnectorSpec.ListenPort != 0 {
		port = resource.Spec.ConnectorSpec.ListenPort
	}
	if len(os.Getenv(envProbeImage)) > 0 {
		probeImage = os.Getenv(envProbeImage)
	}
	if len(resource.Spec.ProbeSpec.Image) > 0 {
		probeImage = resource.Spec.ProbeSpec.Image
	}
	if len(resource.Spec.ConnectorSpec.Image) > 0 {
		MSImage = resource.Spec.ConnectorSpec.Image
	}
	args := []string{"--crd-namespace",
		resource.Namespace, "--crd-name", resource.Name, "--crd-kind", resource.Kind,
		"--crd-uid", string(resource.UID), "--patcher-cfg-path", defaultPatcherPath + "/" + defaultPathcerFileName}
	args = append(args, resource.Spec.ProbeSpec.Args...)
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
							Name:    "probe",
							Image:   probeImage,
							Args:    args,
							Resources: corev1.ResourceRequirements{
								Limits: corev1.ResourceList{
									corev1.ResourceCPU: k8sresource.MustParse("500m"),
									corev1.ResourceMemory: k8sresource.MustParse("2Gi"),
								},
							},
							VolumeMounts: []corev1.VolumeMount{
								{
									Name: "patcherconfig",
									MountPath: defaultPatcherPath,
								},
							},
						},
						{
							Name:    "connector",
							Image:   MSImage,
							Resources: corev1.ResourceRequirements{
								Limits: corev1.ResourceList{
									corev1.ResourceCPU: k8sresource.MustParse("500m"),
									corev1.ResourceMemory: k8sresource.MustParse("2Gi"),
								},
							},
							Env: []corev1.EnvVar{
								{
									Name: EnvMQBrokerHost,
									Value: defaultMQBrokerHost,
								},
								{
									Name: EnvSpringBootAdminHost,
									Value: defaultSpringBootAdminHost,
								},
							},
							Ports: []corev1.ContainerPort{
								{
									Name: defaultMSPortName,
									ContainerPort: port,
								},
							},
						},
					},
					Volumes: []corev1.Volume{
						{
							Name: "patcherconfig",
							VolumeSource: corev1.VolumeSource{
								ConfigMap: &corev1.ConfigMapVolumeSource{
									LocalObjectReference: corev1.LocalObjectReference{
										Name: resource.Name + "-cm",
									},
									Items: []corev1.KeyToPath{
										{
											Key: "patcherconfig",
											Path: defaultPathcerFileName,
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
}

func newServiceForMS(resource *fusionappv1alpha1.Resource) *corev1.Service {
	port := defaultMSPort
	if resource.Spec.ConnectorSpec.ListenPort != 0 {
		port = resource.Spec.ConnectorSpec.ListenPort
	}
	service := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name: resource.Name + "ms-svc",
			Namespace: resource.Namespace,
			Labels: DefaultLabels(resource),
		},
		Spec:       corev1.ServiceSpec{
			Type: corev1.ServiceTypeClusterIP,
			Selector: DefaultLabels(resource),
			Ports: []corev1.ServicePort{
				{
					Name: defaultMSPortName,
					Port: port,
				},
			},
		},
	}
	return service
}

func GatewayMappingGVK() schema.GroupVersionKind {
	return schema.GroupVersionKind{
		Group:   "getambassador.io",
		Kind:    "Mapping",
		Version: "v1",
	}
}

func (r *ReconcileResource) newGatewayMappingForService(resource *fusionappv1alpha1.Resource) (*unstructured.Unstructured, error) {
	port := defaultMSPort
	if resource.Spec.ConnectorSpec.ListenPort != 0 {
		port = resource.Spec.ConnectorSpec.ListenPort
	}
	mapping := &unstructured.Unstructured{}
	mapping.Object = map[string]interface{}{
		"metadata": map[string]interface{}{
			"name":      resource.Name + "-mapping",
			"namespace": resource.Namespace,
			"labels":    DefaultLabels(resource),
		},
		"spec": map[string]interface{}{
			"prefix":        "/" + resource.Name + "/",
			"service":       fmt.Sprintf("%s.%s.svc.cluster.local:%d", resource.Name + "ms-svc", resource.Namespace, port),
		},
	}
	mapping.SetGroupVersionKind(GatewayMappingGVK())
	if err := controllerutil.SetControllerReference(resource, mapping, r.scheme); err != nil {
		return nil, err
	}
	return mapping, nil
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
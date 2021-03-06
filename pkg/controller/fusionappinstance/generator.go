package fusionappinstance

import (
	"github.com/fusion-app/fusion-app/pkg/apis/fusionapp/v1alpha1"
	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	k8sresource "k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"os"
)

const (
	defaultProbeImage      = "registry.cn-shanghai.aliyuncs.com/fusion-app/http-prober:instance-prober.201911292308"

	topic                  = "resource-event-source"

	defaultMqAddress       = "221.228.66.83:30595"    // "114.212.87.225:32015"

	EnvMqAdress            = "MQ_ADRESS"

	envProbeImage          = "APPINSTANCE_PROBE_IMAGE"

	EnvAppInstanceHttpURL  = "APPINSTANCE_HTTP_URL"

	defaultHttpURL         = "https://www.cpss2019.fun:5001/get_app_instance_action_state_and_resource_by_uid"
)

// newDeployForProbe returns a probe deployment with the same name/namespace as the cr
func newDeployForProbe(appInstance *v1alpha1.FusionAppInstance) *appsv1.Deployment {
	probeImage := defaultProbeImage
	if len(os.Getenv(envProbeImage)) > 0 {
		probeImage = os.Getenv(envProbeImage)
	}
	if len(appInstance.Spec.ProbeImage) > 0 {
		probeImage = appInstance.Spec.ProbeImage
	}
	mqAddress := os.Getenv(EnvMqAdress)
	if len(mqAddress) == 0 {
		mqAddress = defaultMqAddress
	}
	httpURL := os.Getenv(EnvAppInstanceHttpURL)
	if len(httpURL) == 0 {
		httpURL = defaultHttpURL
	}
	args := []string{"--mq-address", mqAddress, "--mq-topic", topic, "--crd-namespace",
		appInstance.Namespace, "--crd-name", appInstance.Name, "--crd-kind", appInstance.Kind,
		"--crd-uid", string(appInstance.UID), "--http-url", httpURL}
	args = append(args, appInstance.Spec.ProbeArgs...)
	return &appsv1.Deployment {
		ObjectMeta: metav1.ObjectMeta{
			Name:      appInstance.Name + "-probe-deploy",
			Namespace: appInstance.Namespace,
			Labels:    DefaultLabels(appInstance),
		},
		Spec: appsv1.DeploymentSpec{
			Selector: metav1.SetAsLabelSelector(DefaultLabels(appInstance)),
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels:    DefaultLabels(appInstance),
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{
						{
							Name:    appInstance.Name,
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

package internal

import (
	"os"

	"k8s.io/api/core/v1"
)

const (
	defaultStorageClass = "nfs-client"
	defaultStorageSize  = "1Gi"
)

func GetStorageClassName() string {
	storageClass := os.Getenv("STORAGE_CLASS")
	if len(storageClass) == 0 {
		return defaultStorageClass
	}
	return storageClass
}

type Volume struct {
	Name      string
	MountPath string
	ReadOnly  bool
}

func (v Volume) AddToPodSpec(spec *v1.PodSpec) {
	spec.Volumes = append(spec.Volumes, v1.Volume{
		Name: v.Name,
		VolumeSource: v1.VolumeSource{
			PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
				ClaimName: v.Name,
			},
		},
	})
	if len(spec.Containers) > 0 {
		spec.Containers[0].VolumeMounts = append(spec.Containers[0].VolumeMounts, v1.VolumeMount{
			MountPath: v.MountPath,
			Name:      v.Name,
			ReadOnly:  v.ReadOnly,
		})
	}
}

func (v Volume) AddToPod(pod *v1.Pod) {
	v.AddToPodSpec(&pod.Spec)
}

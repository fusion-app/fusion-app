package types

import (
	"github.com/fusion-app/fusion-app/pkg/apis/fusionapp/v1alpha1"
)

func V1alpha1ResourceToResource(rs *v1alpha1.Resource) *Resource {
	resource := new(Resource)
	resource.UID = string(rs.UID)
	resource.Namespace = rs.Namespace
	resource.Kind = string(rs.Spec.ResourceKind)
	resource.ProbeEnabled = rs.Spec.ProbeEnabled
	resource.Phase = string(rs.Status.ProbePhase)
	resource.Bound = rs.Status.Bound
	resource.Name = rs.Name
	resource.AliasName =  rs.Spec.AliasName
	resource.AccessMode = string(rs.Spec.AccessMode)
	if rs.Spec.Labels != nil {
		in, out := &rs.Spec.Labels, &resource.Labels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if rs.Spec.Operation != nil {
		in, out := &rs.Spec.Operation, &resource.Operation
		*out = make([]v1alpha1.ResourceOperationSpec, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if rs.Spec.ProbeArgs != nil {
		in, out := &rs.Spec.ProbeArgs, &resource.ProbeArgs
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	resource.Icon = rs.Spec.Icon
	if rs.Spec.Description != nil {
		in, out := &rs.Spec.Description, &resource.Description
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	return resource
}

func ResourceToV1alpha1Resource(resource *Resource) *v1alpha1.Resource {
	rs := new(v1alpha1.Resource)
	rs.Namespace = resource.Namespace
	rs.Name = resource.Name
	rs.Spec.ResourceKind = v1alpha1.ResourceKind(resource.Kind)
	rs.Spec.ProbeEnabled = resource.ProbeEnabled
	rs.Status.Bound = resource.Bound
	rs.Spec.AccessMode = v1alpha1.ResourceAccessMode(resource.AccessMode)
	rs.Spec.ProbeImage = resource.ProbeImage
	rs.Spec.AliasName = resource.AliasName
	if resource.Labels != nil {
		in, out := &resource.Labels, &rs.Spec.Labels
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if resource.Operation != nil {
		in, out := &resource.Operation, &rs.Spec.Operation
		*out = make([]v1alpha1.ResourceOperationSpec, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if resource.ProbeArgs != nil {
		in, out := &resource.ProbeArgs, &rs.Spec.ProbeArgs
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	rs.Spec.Icon = resource.Icon
	if resource.Description != nil {
		in, out := &resource.Description, &rs.Spec.Description
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	return rs
}

func UpdateResourceWithResourceSpec(resource *v1alpha1.Resource, spec *ResourceSpec) {
	if spec.Labels != nil {
		if resource.Spec.Labels == nil {
			resource.Spec.Labels = make(map[string]string)
		}
		for key, val := range spec.Labels {
			resource.Spec.Labels[key] = val
		}
	}
	if spec.ProbeArgs != nil {
		in, out := &spec.ProbeArgs, &resource.Spec.ProbeArgs
		*out = make([]string, len(*in))
		copy(*out, *in)
	}
	if len(spec.AliasName) > 0 {
		resource.Spec.AliasName = spec.AliasName
	}
	if spec.Description != nil {
		in, out := &spec.Description, &resource.Spec.Description
		*out = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	if len(spec.Icon) > 0 {
		resource.Spec.Icon = spec.Icon
	}
	if len(spec.ProbeImage) > 0 {
		resource.Spec.ProbeImage = spec.ProbeImage
	}
	if len(spec.AccessMode) > 0 {
		resource.Spec.AccessMode = spec.AccessMode
	}
	if spec.Operation != nil && len(spec.Operation) > 0 {
		in, out := &spec.Operation, &resource.Spec.Operation
		*out = make([]v1alpha1.ResourceOperationSpec, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	resource.Spec.ProbeEnabled = spec.ProbeEnabled
}

func V1alpha1AppInstanceToAppInstance(fusionAppInstance *v1alpha1.FusionAppInstance) *AppInstance {
	appInstance := new(AppInstance)
	appInstance.UID = string(fusionAppInstance.UID)
	appInstance.Namespace = fusionAppInstance.Namespace
	appInstance.Name = fusionAppInstance.Name
	appInstance.RefApp.Name = fusionAppInstance.Spec.RefApp.Name
	appInstance.RefApp.UID = fusionAppInstance.Spec.RefApp.UID
	if fusionAppInstance.Spec.RefResource != nil {
		appInstance.RefResource = make([]AppRefResource, len(fusionAppInstance.Spec.RefResource))
		for i, refResource := range fusionAppInstance.Spec.RefResource {
			appInstance.RefResource[i].UID = refResource.UID
			appInstance.RefResource[i].Name = refResource.Name
			appInstance.RefResource[i].Kind = refResource.Kind
			appInstance.RefResource[i].Namespace = refResource.Namespace
		}
	}
	appInstance.Status = string(fusionAppInstance.Status.Phase)
	if fusionAppInstance.Status.StartTime != nil {
		appInstance.StartTime = fusionAppInstance.Status.StartTime.String()
	}
	if fusionAppInstance.Status.UpdateTime != nil {
		appInstance.UpdateTime = fusionAppInstance.Status.UpdateTime.String()
	}
	if fusionAppInstance.Status.EndTime != nil {
		appInstance.EndTime = fusionAppInstance.Status.EndTime.String()
	}
	return appInstance
}
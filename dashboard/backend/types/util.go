package types

import (
	"github.com/fusion-app/fusion-app/pkg/apis/fusionapp/v1alpha1"
	"math/rand"
	"strconv"
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
	if rs.Labels == nil {
		rs.Labels = make(map[string]string)
	}
	if longitudeString, ok := rs.Labels["longitude"]; ok {
		longitude, err := strconv.ParseFloat(longitudeString, 64)
		if err != nil {
			resource.Position.Longitude = rand.Float64()*(RightBound - LeftBound) + LeftBound
		} else {
			resource.Position.Longitude = longitude
		}
	} else {
		resource.Position.Longitude = rand.Float64()*(RightBound - LeftBound) + LeftBound
		rs.Labels["longitude"] = strconv.FormatFloat(resource.Position.Longitude, 'f', -1, 64)
	}
	if latitudeString, ok := rs.Labels["latitude"]; ok {
		latitude, err := strconv.ParseFloat(latitudeString, 64)
		if err != nil {
			resource.Position.Latitude = rand.Float64()*(UpBound - DownBound) + DownBound
		} else {
			resource.Position.Latitude = latitude
		}
	} else {
		resource.Position.Latitude = rand.Float64()*(UpBound - DownBound) + DownBound
		rs.Labels["latitude"] = strconv.FormatFloat(resource.Position.Latitude, 'f', -1, 64)
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

func V1alpha1AppToApp(fusionApp *v1alpha1.FusionApp) *App {
	app := new(App)
	app.UID = string(fusionApp.UID)
	app.Name = fusionApp.Name
	if fusionApp.Spec.ResourceClaim != nil {
		in, out := &fusionApp.Spec.ResourceClaim, &app.ResourceClaim
		*out = make([]v1alpha1.ResourceClaimSpec, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	app.Icon = fusionApp.Spec.Icon
	app.AliasName = fusionApp.Spec.AliasName
	if fusionApp.Spec.Description != nil {
		in, out := &fusionApp.Spec.Description, &app.Description
		*out  = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	return app
}

func AppToV1alpha1App(app *App) *v1alpha1.FusionApp {
	fusionApp := new(v1alpha1.FusionApp)
	fusionApp.Name = app.Name
	if app.ResourceClaim != nil {
		in, out := &app.ResourceClaim, &fusionApp.Spec.ResourceClaim
		*out = make([]v1alpha1.ResourceClaimSpec, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	fusionApp.Spec.Icon = app.Icon
	fusionApp.Spec.AliasName = app.AliasName
	if app.Description != nil {
		in, out := &app.Description, &fusionApp.Spec.Description
		*out  = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
	return fusionApp
}

func UpdateAppWithAppSpec(fusionApp *v1alpha1.FusionApp, app *App)  {
	if app.ResourceClaim != nil {
		in, out := &app.ResourceClaim, &fusionApp.Spec.ResourceClaim
		*out = make([]v1alpha1.ResourceClaimSpec, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if len(app.Icon) != 0 {
		fusionApp.Spec.Icon = app.Icon
	}
	if len(app.AliasName) != 0 {
		fusionApp.Spec.AliasName = app.AliasName
	}
	if app.Description != nil {
		in, out := &app.Description, &fusionApp.Spec.Description
		*out  = make(map[string]string, len(*in))
		for key, val := range *in {
			(*out)[key] = val
		}
	}
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
	if fusionAppInstance.Status.ActionStatus != nil {
		in, out := &fusionAppInstance.Status.ActionStatus, &appInstance.Status
		*out = make([]v1alpha1.Action, len(*in))
		for i := range *in {
			(*in)[i].DeepCopyInto(&(*out)[i])
		}
	}
	if fusionAppInstance.Status.CreateTime != nil {
		appInstance.CreateTime = fusionAppInstance.Status.CreateTime.String()
	}
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
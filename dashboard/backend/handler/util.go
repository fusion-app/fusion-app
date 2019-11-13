package handler

import "github.com/fusion-app/fusion-app/pkg/apis/fusionapp/v1alpha1"

func v1alpha1resourceToResource(rs *v1alpha1.Resource) Resource {
	var resource Resource
	resource.UID = string(rs.UID)
	resource.Namespace = rs.Namespace
	resource.Kind = string(rs.Spec.ResourceKind)
	resource.Phase = string(rs.Status.ProbePhase)
	resource.Bound = rs.Status.Bound
	resource.Name = rs.Name
	resource.AccessMode = string(rs.Spec.AccessMode)
	resource.Labels = rs.Labels
	resource.Operation = rs.Spec.Operation
	resource.Icon = rs.Spec.Icon
	resource.Description = rs.Spec.Description
	return resource
}

func resourceToV1alpha1Resource(resource *Resource) v1alpha1.Resource {
	rs := new(v1alpha1.Resource)
	rs.Namespace = resource.Namespace
	rs.Name = resource.Name
	rs.Spec.ResourceKind = v1alpha1.ResourceKind(resource.Kind)
	rs.Status.Bound = resource.Bound
	rs.Spec.AccessMode = v1alpha1.ResourceAccessMode(resource.AccessMode)
	rs.Labels = resource.Labels
	rs.Spec.Operation = resource.Operation
	rs.Spec.Icon = resource.Icon
	rs.Spec.Description = resource.Description
	return *rs
}

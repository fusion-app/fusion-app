package util

import (
	"context"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/types"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func CreateIfNotExists(c client.Client, ifExistKey client.ObjectKey, obj runtime.Object, gvk schema.GroupVersionKind) error {
	found := &unstructured.Unstructured{}
	found.SetGroupVersionKind(gvk)
	err := c.Get(context.Background(), ifExistKey, found)
	if err != nil && errors.IsNotFound(err) {
		return c.Create(context.TODO(), obj)
	} else if err != nil {
		return err
	} else {
		return nil
	}
}

func CreateIfNotExistsMapping(c client.Client, mapping *unstructured.Unstructured) error {
	return CreateIfNotExists(c, types.NamespacedName{
		Namespace: mapping.GetNamespace(),
		Name:      mapping.GetName(),
	}, mapping, schema.GroupVersionKind{
		Group:   "getambassador.io",
		Kind:    "Mapping",
		Version: "v1",
	})
}

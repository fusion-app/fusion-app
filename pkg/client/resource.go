package client

import (
	"github.com/fusion-app/fusion-app/pkg/apis/fusionapp/v1alpha1"
	"github.com/kubeless/kubeless/pkg/client/clientset/versioned/scheme"
	"k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/watch"
	"k8s.io/client-go/rest"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

type ResourceInterface interface {
	Create(*v1alpha1.Resource) (*v1alpha1.Resource, error)
	Update(*v1alpha1.Resource) (*v1alpha1.Resource, error)
	Delete(name string, options *v1.DeleteOptions) error
	DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error
	Get(name string, options v1.GetOptions) (*v1alpha1.Resource, error)
	List(opts v1.ListOptions) (*v1alpha1.ResourceList, error)
	Watch(opts v1.ListOptions) (watch.Interface, error)
	Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Resource, err error)
}

// resources implements ResourceInterface
type resources struct {
	client rest.Interface
	ns     string
}

func NewResources(namespace string) (*resources, error) {
	kubeConfig, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	client, err := rest.RESTClientFor(kubeConfig)
	if err != nil {
		return nil, err
	}
	return &resources{
		client: client,
		ns:     namespace,
	}, nil
}

// Get takes name of the function, and returns the corresponding function object, and an error if there is any.
func (c *resources) Get(name string, options v1.GetOptions) (result *v1alpha1.Resource, err error) {
	result = &v1alpha1.Resource{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("resources").
		Name(name).
		VersionedParams(&options, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// List takes label and field selectors, and returns the list of Resources that match those selectors.
func (c *resources) List(opts v1.ListOptions) (result *v1alpha1.ResourceList, err error) {
	result = &v1alpha1.ResourceList{}
	err = c.client.Get().
		Namespace(c.ns).
		Resource("resources").
		VersionedParams(&opts, scheme.ParameterCodec).
		Do().
		Into(result)
	return
}

// Watch returns a watch.Interface that watches the requested resources.
func (c *resources) Watch(opts v1.ListOptions) (watch.Interface, error) {
	opts.Watch = true
	return c.client.Get().
		Namespace(c.ns).
		Resource("resources").
		VersionedParams(&opts, scheme.ParameterCodec).
		Watch()
}

// Create takes the representation of a function and creates it.  Returns the server's representation of the function, and an error, if there is any.
func (c *resources) Create(function *v1alpha1.Resource) (result *v1alpha1.Resource, err error) {
	result = &v1alpha1.Resource{}
	err = c.client.Post().
		Namespace(c.ns).
		Resource("resources").
		Body(function).
		Do().
		Into(result)
	return
}

// Update takes the representation of a function and updates it. Returns the server's representation of the function, and an error, if there is any.
func (c *resources) Update(function *v1alpha1.Resource) (result *v1alpha1.Resource, err error) {
	result = &v1alpha1.Resource{}
	err = c.client.Put().
		Namespace(c.ns).
		Resource("resources").
		Name(function.Name).
		Body(function).
		Do().
		Into(result)
	return
}

// Delete takes name of the function and deletes it. Returns an error if one occurs.
func (c *resources) Delete(name string, options *v1.DeleteOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("resources").
		Name(name).
		Body(options).
		Do().
		Error()
}

// DeleteCollection deletes a collection of objects.
func (c *resources) DeleteCollection(options *v1.DeleteOptions, listOptions v1.ListOptions) error {
	return c.client.Delete().
		Namespace(c.ns).
		Resource("resources").
		VersionedParams(&listOptions, scheme.ParameterCodec).
		Body(options).
		Do().
		Error()
}

// Patch applies the patch and returns the patched function.
func (c *resources) Patch(name string, pt types.PatchType, data []byte, subresources ...string) (result *v1alpha1.Resource, err error) {
	result = &v1alpha1.Resource{}
	err = c.client.Patch(pt).
		Namespace(c.ns).
		Resource("resources").
		SubResource(subresources...).
		Name(name).
		Body(data).
		Do().
		Into(result)
	return
}

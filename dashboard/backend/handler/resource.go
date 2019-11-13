package handler

import (
	"context"
	"encoding/json"
	"fmt"
	resourcev1alpha1 "github.com/fusion-app/fusion-app/pkg/apis/fusionapp/v1alpha1"
	resourcecontroller "github.com/fusion-app/fusion-app/pkg/controller/resource"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/api/errors"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (handler *APIHandler) ListResourcesWithKind(w http.ResponseWriter, r *http.Request) {
	resourceAPIQueryBody := new(ResourceAPIQueryBody)
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &resourceAPIQueryBody); err != nil {
		if err := json.NewEncoder(w).Encode(err); err != nil {
			responseJSON(Message{err.Error()}, w, http.StatusUnprocessableEntity)
		}
	}
	rsl := &resourcev1alpha1.ResourceList{}
	err = handler.client.List(context.TODO(), &client.ListOptions{}, rsl)
	if err != nil {
		log.Warningf("failed to list resources: %v", err)
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	}
	var resources []Resource
	for _, resource := range rsl.Items {
		if (len(resourceAPIQueryBody.Kind) == 0 || string(resource.Spec.ResourceKind) == resourceAPIQueryBody.Kind) &&
			(len(resourceAPIQueryBody.Phase) == 0 || string(resource.Status.Phase) == resourceAPIQueryBody.Phase) {
			resources = append(resources, *v1alpha1resourceToResource(&resource))
		}
	}
	responseJSON(resources, w, http.StatusOK)
}

func (handler *APIHandler) CreateResource(w http.ResponseWriter, r *http.Request) {
	resource := new(Resource)
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	}
	defer r.Body.Close()
	if err := json.Unmarshal(body, &resource); err != nil {
		if err := json.NewEncoder(w).Encode(err); err != nil {
			responseJSON(Message{err.Error()}, w, http.StatusUnprocessableEntity)
		}
	}
	if len(resource.Namespace) == 0 {
		resource.Namespace = handler.resourcesNamespace
	}
	rs := resourceToV1alpha1Resource(resource)
	resourcecontroller.AddKindLabel(rs)

	err = handler.client.Create(context.TODO(), rs)
	if err != nil {
		log.Warningf("Failed to create resource %v: %v", resource.Name, err)
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	} else {
		responseJSON(resource, w, http.StatusCreated)
	}
}

func (handler *APIHandler) UpdateResource(w http.ResponseWriter, r *http.Request) {
	resourceAPIPutBody := new(ResourceAPIPutBody)
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &resourceAPIPutBody); err != nil {
		if err := json.NewEncoder(w).Encode(err); err != nil {
			responseJSON(Message{err.Error()}, w, http.StatusUnprocessableEntity)
		}
	}
	name := resourceAPIPutBody.AppRefResource.Name
	namespace := resourceAPIPutBody.AppRefResource.Namespace
	if len(namespace) == 0 {
		namespace = handler.resourcesNamespace
	}
	resource := new(resourcev1alpha1.Resource)
	err = handler.client.Get(context.TODO(), client.ObjectKey{Namespace: namespace,
		Name: name}, resource)
	if errors.IsNotFound(err) {
		err := fmt.Errorf("resource \"%s\" not exists", )
		responseJSON(Message{err.Error()}, w, http.StatusNotFound)
		return
	} else if err != nil {
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
		return
	}

	if string(resource.Spec.ResourceKind) != resourceAPIPutBody.AppRefResource.Kind {
		err := fmt.Errorf("resource kind is not correct")
		responseJSON(Message{err.Error()}, w, http.StatusBadRequest)
		return
	}
	updateResourceWithResourceSpec(resource, &resourceAPIPutBody.ResourceSpec)
	err = handler.client.Update(context.TODO(), resource)
	if err != nil {
		log.Warningf("Failed to create workspace %v: %v", resource.Name, err)
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	} else {
		responseJSON(resource, w, http.StatusOK)
	}
}

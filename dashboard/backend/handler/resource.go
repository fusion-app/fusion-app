package handler

import (
	"context"
	"encoding/json"
	"fmt"
	resourcev1alpha1 "github.com/fusion-app/fusion-app/pkg/apis/fusionapp/v1alpha1"
	resourcecontroller "github.com/fusion-app/fusion-app/pkg/controller/resource"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/api/errors"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (handler *APIHandler) ListResources(w http.ResponseWriter, r *http.Request) {
	rsl := &resourcev1alpha1.ResourceList{}

	err := handler.client.List(context.TODO(), &client.ListOptions{}, rsl)
	var resources []Resource
	for _, resource := range rsl.Items {
		resources = append(resources, v1alpha1resourceToResource(&resource))
	}
	if err != nil {
		log.Warningf("failed to list resources: %v", err)
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	} else {
		responseJSON(ResourceList{Resources: resources}, w, http.StatusOK)
	}
}

func (handler *APIHandler) ListResourcesWithKind(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	kind := vars["kind"]
	rsl := &resourcev1alpha1.ResourceList{}
	err := handler.client.List(context.TODO(), &client.ListOptions{}, rsl)
	var resources []Resource
	for _, resource := range rsl.Items {
		if string(resource.Spec.ResourceKind) == kind {
			resources = append(resources, v1alpha1resourceToResource(&resource))
		}
	}
	if err != nil {
		log.Warningf("failed to list resources: %v", err)
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	} else {
		responseJSON(ResourceList{Resources: resources}, w, http.StatusOK)
	}
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
	resourcecontroller.AddKindLabel(&rs)

	err = handler.client.Create(context.TODO(), &rs)
	if err != nil {
		log.Warningf("Failed to create resource %v: %v", resource.Name, err)
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	} else {
		responseJSON(resource, w, http.StatusCreated)
	}
}

func (handler *APIHandler) UpdateResource(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	kind := vars["kind"]
	name := vars["resource"]

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

	oldResource := new(resourcev1alpha1.Resource)
	err = handler.client.Get(context.TODO(), client.ObjectKey{Namespace: resource.Namespace, Name: name}, oldResource)
	if errors.IsNotFound(err) {
		err := fmt.Errorf("resource \"%s\" not exists", name)
		responseJSON(Message{err.Error()}, w, http.StatusNotFound)
		return
	} else if err != nil {
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
		return
	}

	if resource.Name != name {
		err := fmt.Errorf("resource name in path is not the same as that in json")
		responseJSON(Message{err.Error()}, w, http.StatusBadRequest)
		return
	}

	if resource.Kind != kind || string(oldResource.Spec.ResourceKind) != kind  {
		err := fmt.Errorf("resource kind in path is not the same as that in json")
		responseJSON(Message{err.Error()}, w, http.StatusBadRequest)
		return
	}
	rs := resourceToV1alpha1Resource(resource)
	err = handler.client.Update(context.TODO(), &rs)
	if err != nil {
		log.Warningf("Failed to create workspace %v: %v", resource.Name, err)
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	} else {
		responseJSON(resource, w, http.StatusOK)
	}
}

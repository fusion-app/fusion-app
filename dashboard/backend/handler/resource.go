package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fusion-app/fusion-app/dashboard/backend/types"
	fusionappv1alpha1 "github.com/fusion-app/fusion-app/pkg/apis/fusionapp/v1alpha1"
	resourcecontroller "github.com/fusion-app/fusion-app/pkg/controller/resource"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"net/http"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (handler *APIHandler) ListResourcesWithKind(w http.ResponseWriter, r *http.Request) {
	resourceAPIQueryBody := new(types.ResourceAPIQueryBody)
	gatewayUrl := os.Getenv("GATEWAY")
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	defer r.Body.Close()
	if err != nil {
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
		return
	}
	if len(body) != 0 {
		if err := json.Unmarshal(body, &resourceAPIQueryBody); err != nil {
			if nerr := json.NewEncoder(w).Encode(err); nerr != nil {
				responseJSON(Message{nerr.Error()}, w, http.StatusUnprocessableEntity)
			} else {
				responseJSON(Message{err.Error()}, w, http.StatusBadRequest)
			}
			return
		}
	}
	if len(resourceAPIQueryBody.RefResource.Name) > 0  {
		resource := &fusionappv1alpha1.Resource{}
		name := resourceAPIQueryBody.RefResource.Name
		namespace := resourceAPIQueryBody.RefResource.Namespace
		if len(namespace) == 0 {
			namespace = handler.resourcesNamespace
		}
		err = handler.client.Get(context.TODO(), client.ObjectKey{Namespace: namespace,
			Name: name}, resource)
		if errors.IsNotFound(err) {
			err := fmt.Errorf("resource \"%s\" not exists", name)
			responseJSON(Message{err.Error()}, w, http.StatusNotFound)
			return
		} else if err != nil {
			responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
			return
		}
		rs, modified := types.V1alpha1ResourceToResource(resource)
		rs.Gateway = gatewayUrl + resource.Name
		resources := []types.Resource{*rs}
		responseJSON(resources, w, http.StatusOK)
		if modified {
			_ = handler.client.Update(context.TODO(), resource)
		}
	} else {
		rsl := &fusionappv1alpha1.ResourceList{}
		err = handler.client.List(context.TODO(), &client.ListOptions{}, rsl)
		if err != nil {
			log.Warningf("failed to list resources", err)
			responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
			return
		}
		resources := make([]types.Resource, 0)
		rss := make([]fusionappv1alpha1.Resource, 0)
		if resourceAPIQueryBody.LabelSelector != nil && len(resourceAPIQueryBody.LabelSelector) > 0 {
			mp := make(labels.Set)
			for _, selector := range resourceAPIQueryBody.LabelSelector {
				mp[selector.Key] = selector.Value
			}
			labelSelector := labels.SelectorFromSet(mp)
			for _, item := range rsl.Items {
				if labelSelector.Matches(labels.Set(item.Spec.Labels)) {
					rs, modified := types.V1alpha1ResourceToResource(&item)
					rs.Gateway = gatewayUrl + item.Name
					resources = append(resources, *rs)
					if modified {
						rss = append(rss, item)
					}
				}
			}
		} else {
			for _, resource := range rsl.Items {
				if (len(resourceAPIQueryBody.Kind) == 0 || string(resource.Spec.ResourceKind) == resourceAPIQueryBody.Kind) &&
					(len(resourceAPIQueryBody.Phase) == 0 || string(resource.Status.ProbePhase) == resourceAPIQueryBody.Phase) &&
					(len(resourceAPIQueryBody.Bound) == 0 || (resource.Status.Bound == true && resourceAPIQueryBody.Bound == "true") || (resource.Status.Bound == false && resourceAPIQueryBody.Bound == "false")){
					rs, modified := types.V1alpha1ResourceToResource(&resource)
					rs.Gateway = gatewayUrl + resource.Name
					resources = append(resources, *rs)
					if modified {
						rss = append(rss, resource)
					}
				}
			}
		}
		responseJSON(resources, w, http.StatusOK)
		for _, rs := range rss {
			_ = handler.client.Update(context.TODO(), &rs)
		}
	}
}

func (handler *APIHandler) CreateResource(w http.ResponseWriter, r *http.Request) {
	resourceAPICreateBody := new(types.ResourceAPICreateBody)
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	defer r.Body.Close()
	if err != nil {
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
		return
	}
	if err := json.Unmarshal(body, &resourceAPICreateBody); err != nil {
		if nerr := json.NewEncoder(w).Encode(err); nerr != nil {
			responseJSON(Message{nerr.Error()}, w, http.StatusUnprocessableEntity)
		} else {
			responseJSON(Message{err.Error()}, w, http.StatusBadRequest)
		}
		return
	}
	resource := resourceAPICreateBody.ResourceSpec
	if len(resource.Namespace) == 0 {
		resource.Namespace = handler.resourcesNamespace
	}
	rs := types.ResourceToV1alpha1Resource(&resource)
	resourcecontroller.AddKindLabel(rs)

	err = handler.client.Create(context.TODO(), rs)
	if err != nil {
		log.Warningf("Failed to create resource %v: %v", resource.Name, err)
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	} else {
		responseJSON(resource, w, http.StatusOK)
	}
}

func (handler *APIHandler) UpdateResource(w http.ResponseWriter, r *http.Request) {
	resourceAPIPutBody := new(types.ResourceAPIPutBody)
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	defer r.Body.Close()
	if err != nil {
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(body, &resourceAPIPutBody); err != nil {
		if nerr := json.NewEncoder(w).Encode(err); nerr != nil {
			responseJSON(Message{nerr.Error()}, w, http.StatusUnprocessableEntity)
		} else {
			responseJSON(Message{err.Error()}, w, http.StatusBadRequest)
		}
		return
	}
	name := resourceAPIPutBody.AppRefResource.Name
	namespace := resourceAPIPutBody.AppRefResource.Namespace
	if len(namespace) == 0 {
		namespace = handler.resourcesNamespace
	}
	resource := new(fusionappv1alpha1.Resource)
	err = handler.client.Get(context.TODO(), client.ObjectKey{Namespace: namespace,
		Name: name}, resource)
	if errors.IsNotFound(err) {
		err := fmt.Errorf("resource \"%s\" not exists", name)
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
	types.UpdateResourceWithResourceSpec(resource, &resourceAPIPutBody.ResourceSpec)
	err = handler.client.Update(context.TODO(), resource)
	if err != nil {
		log.Warningf("Failed to Update resource %v: %v", resource.Name, err)
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	} else {
		responseJSON(resource, w, http.StatusOK)
	}
}

func (handler *APIHandler) BindResource(w http.ResponseWriter, r *http.Request)  {
	handler.handleBind(w, r, true)
}

func (handler *APIHandler) UnBindResource(w http.ResponseWriter, r *http.Request)  {
	handler.handleBind(w, r, false)
}

func (handler *APIHandler) handleBind(w http.ResponseWriter, r *http.Request, bind bool) {
	resourceAPIBindBody := new(types.ResourceAPIBindBody)
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	defer r.Body.Close()
	if err != nil {
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(body, &resourceAPIBindBody); err != nil {
		if nerr := json.NewEncoder(w).Encode(err); nerr != nil {
			responseJSON(Message{nerr.Error()}, w, http.StatusUnprocessableEntity)
		} else {
			responseJSON(Message{err.Error()}, w, http.StatusBadRequest)
		}
		return
	}
	name := resourceAPIBindBody.RefResource.Name
	namespace := resourceAPIBindBody.RefResource.Namespace
	if len(namespace) == 0 {
		namespace = handler.resourcesNamespace
	}
	resource := new(fusionappv1alpha1.Resource)
	err = handler.client.Get(context.TODO(), client.ObjectKey{Namespace: namespace,
		Name: name}, resource)
	if errors.IsNotFound(err) {
		err := fmt.Errorf("resource \"%s\" not exists", name)
		responseJSON(Message{err.Error()}, w, http.StatusNotFound)
		return
	} else if err != nil {
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
		return
	}
	resource.Status.Bound = bind
	err = handler.client.Update(context.TODO(), resource)
	if err != nil {
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
		return
	}
	responseJSON("ok", w, http.StatusOK)
}
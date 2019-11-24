package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/fusion-app/fusion-app/dashboard/backend/types"
	fusionappv1alpha1 "github.com/fusion-app/fusion-app/pkg/apis/fusionapp/v1alpha1"
	"github.com/fusion-app/fusion-app/pkg/util"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/labels"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sort"
)

func (handler *APIHandler) QueryAppInstance(w http.ResponseWriter, r *http.Request) {
	appInstanceAPIQueryBody := new(types.AppInstanceAPIQueryBody)
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	defer r.Body.Close()
	if err != nil {
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(body, &appInstanceAPIQueryBody); err != nil {
		if nerr := json.NewEncoder(w).Encode(err); nerr != nil {
			responseJSON(Message{nerr.Error()}, w, http.StatusUnprocessableEntity)
		} else {
			responseJSON(Message{err.Error()}, w, http.StatusBadRequest)
		}
		return
	}
	refAppInstance := appInstanceAPIQueryBody.RefAppInstance
	name := refAppInstance.Name
	namespace := refAppInstance.Namespace
	if len(namespace) == 0 {
		namespace = handler.resourcesNamespace
	}
	fusionAppInstance := &fusionappv1alpha1.FusionAppInstance{}
	err = handler.client.Get(context.TODO(), client.ObjectKey{Namespace: namespace,
		Name: name}, fusionAppInstance)
	if errors.IsNotFound(err) {
		err := fmt.Errorf("appinstance \"%s\" not exists", name)
		responseJSON(Message{err.Error()}, w, http.StatusNotFound)
		return
	} else if err != nil {
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
		return
	}
	appInstance := types.V1alpha1AppInstanceToAppInstance(fusionAppInstance)
	responseJSON(appInstance, w, http.StatusOK)
}

func (handler *APIHandler) CreateAppInstance(w http.ResponseWriter, r *http.Request) {
	appInstanceAPICreateBody := new(types.AppInstanceAPICreateBody)
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	defer r.Body.Close()
	if err != nil {
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
		return
	}
	if err := json.Unmarshal(body, &appInstanceAPICreateBody); err != nil {
		if nerr := json.NewEncoder(w).Encode(err); nerr != nil {
			responseJSON(Message{nerr.Error()}, w, http.StatusUnprocessableEntity)
		} else {
			responseJSON(Message{err.Error()}, w, http.StatusBadRequest)
		}
		return
	}
	fusionAppInstance := new(fusionappv1alpha1.FusionAppInstance)
	fusionAppInstance.Namespace = handler.resourcesNamespace
	fusionAppInstance.Name = appInstanceAPICreateBody.RefApp.Name + util.RandRunes(8)
	fusionAppInstance.Spec.RefApp.Name = appInstanceAPICreateBody.RefApp.Name
	app := new(fusionappv1alpha1.FusionApp)
	err = handler.client.Get(context.TODO(), client.ObjectKey{Name: appInstanceAPICreateBody.RefApp.Name, Namespace:
		handler.resourcesNamespace}, app)
	if errors.IsNotFound(err) {
		log.Warningf("failed to get fusionApp %v: %v", appInstanceAPICreateBody.RefApp.Name, err)
		responseJSON(Message{err.Error()}, w, http.StatusBadRequest)
		return
	} else if err != nil {
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
		return
	}
	fusionAppInstance.Spec.RefApp.UID = string(app.UID)

	rsl := &fusionappv1alpha1.ResourceList{}
	err = handler.client.List(context.TODO(), &client.ListOptions{}, rsl)
	if err != nil {
		log.Warningf("failed to list resources", err)
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
		return
	}
	resources := make([]fusionappv1alpha1.Resource, 0)
	for _, item := range rsl.Items {
		if !item.Status.Bound {
			resources = append(resources, item)
		}
	}
	if len(resources) == 0 {
		responseJSON(Message{"No available resources"}, w, http.StatusInternalServerError)
		return
	}
	if app.Spec.ResourceClaim != nil {
		for _, resourceClaim := range app.Spec.ResourceClaim {
			mp := make(labels.Set)
			for _, selector := range resourceClaim.Selector {
				mp[selector.Key] = selector.Value
			}
			if mp.Has("io.fusionapp.smarthome/type") {
				url := "http://nemoworks.info:8080/resource/get/" + mp.Get("io.fusionapp.smarthome/type")
				data, _ := json.Marshal(appInstanceAPICreateBody.UserLabel)
				reader := bytes.NewReader(data)
				request, _ := http.NewRequest("POST", url, reader)
				request.Header.Set("Content-Type", "application/json;charset=UTF-8")
				clt := &http.Client{}
				resp, err := clt.Do(request)
				if err == nil && ( resp.StatusCode == http.StatusOK ){
					respBody := new(types.RespBody)
					body, _ := ioutil.ReadAll(io.LimitReader(resp.Body, 1048576))
					if err := json.Unmarshal(body, &respBody); err == nil {
						_ = resp.Body.Close()
						if respBody.Status == 200 {
							fusionAppInstance.Spec.RefResource = append(fusionAppInstance.Spec.RefResource, fusionappv1alpha1.RefResource{
								Kind: respBody.RespData.SourceDetail.Kind,
								Name: respBody.RespData.SourceDetail.Name,
								Namespace: respBody.RespData.SourceDetail.Namespace,
								UID: respBody.RespData.SourceDetail.UID,
								AliasName: respBody.RespData.SourceDetail.AliasName,
								Icon: respBody.RespData.SourceDetail.Icon,
								Description: respBody.RespData.SourceDetail.Description,
							})
							continue
						}
					} else {
						_ = resp.Body.Close()
					}
				}
			}
			labelSelector := labels.SelectorFromSet(mp)
			var resource *fusionappv1alpha1.Resource
			for _, item := range resources {
				if labelSelector.Matches(labels.Set(item.Spec.Labels)) {
					resource = &item
					break
				}
			}
			if resource == nil {
				responseJSON(Message{"No available resources"}, w, http.StatusInternalServerError)
				return
			} else {
				fusionAppInstance.Spec.RefResource = append(fusionAppInstance.Spec.RefResource, fusionappv1alpha1.RefResource{
					Kind: string(resource.Spec.ResourceKind),
					Name: resource.Name,
					Namespace: resource.Namespace,
					UID: string(resource.UID),
					AliasName: resource.Spec.AliasName,
					Icon: resource.Spec.Icon,
					Description: resource.Spec.Description,
				})
			}
		}
	}
	err = handler.client.Create(context.TODO(), fusionAppInstance)
	if err != nil {
		log.Warningf("failed to create fusionAppInstance %v: %v", fusionAppInstance.Name, err)
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	} else {
		responseJSON(fusionAppInstance, w, http.StatusOK)
	}
}

func (handler *APIHandler) ListAppInstance(w http.ResponseWriter, r *http.Request) {
	appInstanceAPIListBody := new(types.AppInstanceAPIListBody)
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	defer r.Body.Close()
	if err != nil {
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
		return
	}
	asl := new(fusionappv1alpha1.FusionAppInstanceList)
	err = handler.client.List(context.TODO(), &client.ListOptions{}, asl)
	if err != nil {
		log.Warningf("failed to list appInstances", err)
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
		return
	}
	ass := asl.Items
	if len(body) != 0 {
		if err := json.Unmarshal(body, &appInstanceAPIListBody); err != nil {
			if nerr := json.NewEncoder(w).Encode(err); nerr != nil {
				responseJSON(Message{nerr.Error()}, w, http.StatusUnprocessableEntity)
			} else {
				responseJSON(Message{err.Error()}, w, http.StatusBadRequest)
			}
			return
		}
	} else {
		appInstances := make([]types.AppInstance, len(ass))
		for i, instance := range ass {
			appInstances[i] = *types.V1alpha1AppInstanceToAppInstance(&instance)
		}
		responseJSON(appInstances, w, http.StatusOK)
		return
	}
	if appInstanceAPIListBody.SortBy.Field == "startTime" {
		if appInstanceAPIListBody.SortBy.Order {
			sort.SliceStable(ass, func(i, j int) bool { return ass[i].Status.StartTime.Before(ass[j].Status.StartTime)})
		} else {
			sort.SliceStable(ass, func(i, j int) bool { return ass[j].Status.StartTime.Before(ass[i].Status.StartTime)})
		}
	} else if appInstanceAPIListBody.SortBy.Field == "updateTime" {
		if appInstanceAPIListBody.SortBy.Order {
			sort.SliceStable(ass, func(i, j int) bool { return ass[i].Status.UpdateTime.Before(ass[j].Status.UpdateTime)})
		} else {
			sort.SliceStable(ass, func(i, j int) bool { return ass[j].Status.UpdateTime.Before(ass[i].Status.UpdateTime)})
		}
	} else if appInstanceAPIListBody.SortBy.Field == "endTime" {
		if appInstanceAPIListBody.SortBy.Order {
			sort.SliceStable(ass, func(i, j int) bool { return ass[i].Status.EndTime.Before(ass[j].Status.EndTime)})
		} else {
			sort.SliceStable(ass, func(i, j int) bool { return ass[j].Status.EndTime.Before(ass[i].Status.EndTime)})
		}
	} else if appInstanceAPIListBody.SortBy.Field == "name" {
		if appInstanceAPIListBody.SortBy.Order {
			sort.SliceStable(ass, func(i, j int) bool { return ass[i].Name < ass[j].Name })
		} else {
			sort.SliceStable(ass, func(i, j int) bool { return ass[j].Name < ass[i].Name })
		}
	}
	lowerBound := appInstanceAPIListBody.Limit*appInstanceAPIListBody.Page
	upperBound := appInstanceAPIListBody.Limit*(appInstanceAPIListBody.Page+1)
	appInstances := make([]types.AppInstance, 0)
	for i := lowerBound; i < upperBound && i < len(ass); i ++ {
		appInstances = append(appInstances, *types.V1alpha1AppInstanceToAppInstance(&ass[i]))
	}
	responseJSON(appInstances, w, http.StatusOK)
}

func (handler *APIHandler) DeleteAppInstance(w http.ResponseWriter, r *http.Request)  {

}
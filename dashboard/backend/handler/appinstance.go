package handler

import (
	"context"
	"encoding/json"
	"fmt"
	fusionappv1alpha1 "github.com/fusion-app/fusion-app/pkg/apis/fusionapp/v1alpha1"
	"github.com/fusion-app/fusion-app/pkg/util"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/api/errors"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sort"
)

func (handler *APIHandler) QueryAppInstance(w http.ResponseWriter, r *http.Request) {
	appInstanceAPIQueryBody := new(AppInstanceAPIQueryBody)
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
	appInstance := v1alpha1AppInstanceToAppInstance(fusionAppInstance)
	responseJSON(appInstance, w, http.StatusOK)
}

func (handler *APIHandler) CreateAppInstance(w http.ResponseWriter, r *http.Request) {
	appInstanceAPICreateBody := new(AppInstanceAPICreateBody)
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
	fusionAppInstance.Spec.RefApp.UID = appInstanceAPICreateBody.RefApp.UID
	if appInstanceAPICreateBody != nil {
		fusionAppInstance.Spec.RefResource = make([]fusionappv1alpha1.AppRefResource, len(appInstanceAPICreateBody.RefResource))
		for i, refResource := range appInstanceAPICreateBody.RefResource {
			fusionAppInstance.Spec.RefResource[i].UID = refResource.UID
			fusionAppInstance.Spec.RefResource[i].Namespace = refResource.Namespace
			if len(refResource.Namespace) == 0 {
				fusionAppInstance.Spec.RefResource[i].Namespace = handler.resourcesNamespace
			}
			fusionAppInstance.Spec.RefResource[i].Kind = refResource.Kind
			fusionAppInstance.Spec.RefResource[i].Name = refResource.Name
		}
	}
	err = handler.client.Create(context.TODO(), fusionAppInstance)
	if err != nil {
		log.Warningf("failed to create fusionAppInstance %v: %v", fusionAppInstance.Name, err)
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	} else {
		responseJSON(fusionAppInstance, w, http.StatusCreated)
	}
}

func (handler *APIHandler) ListAppInstance(w http.ResponseWriter, r *http.Request) {
	appInstanceAPIListBody := new(AppInstanceAPIListBody)
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	defer r.Body.Close()
	if err != nil {
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	}
	if err := json.Unmarshal(body, &appInstanceAPIListBody); err != nil {
		if nerr := json.NewEncoder(w).Encode(err); nerr != nil {
			responseJSON(Message{nerr.Error()}, w, http.StatusUnprocessableEntity)
		} else {
			responseJSON(Message{err.Error()}, w, http.StatusBadRequest)
		}
		return
	}
	asl := new(fusionappv1alpha1.FusionAppInstanceList)
	err = handler.client.List(context.TODO(), &client.ListOptions{}, asl)
	if err != nil {
		log.Warningf("failed to list appInstances", err)
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	}
	ass := asl.Items
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
	var appInstances []AppInstance
	for i := lowerBound; i < upperBound && i < len(ass); i ++ {
		appInstances = append(appInstances, *v1alpha1AppInstanceToAppInstance(&ass[i]))
	}
	responseJSON(appInstances, w, http.StatusCreated)
}

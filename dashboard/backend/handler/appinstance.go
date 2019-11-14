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
)

func (handler *APIHandler) QueryAppInstance(w http.ResponseWriter, r *http.Request) {
	refAppInstance := new(RefAppInstance)
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	}
	defer r.Body.Close()

	if err := json.Unmarshal(body, &refAppInstance); err != nil {
		if err := json.NewEncoder(w).Encode(err); err != nil {
			responseJSON(Message{err.Error()}, w, http.StatusUnprocessableEntity)
		}
	}
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
	if err != nil {
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	}
	defer r.Body.Close()
	if err := json.Unmarshal(body, &appInstanceAPICreateBody); err != nil {
		if err := json.NewEncoder(w).Encode(err); err != nil {
			responseJSON(Message{err.Error()}, w, http.StatusUnprocessableEntity)
		}
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
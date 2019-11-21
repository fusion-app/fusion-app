package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/fusion-app/fusion-app/dashboard/backend/types"
	fusionappv1alpha1 "github.com/fusion-app/fusion-app/pkg/apis/fusionapp/v1alpha1"
	log "github.com/sirupsen/logrus"
	"io"
	"io/ioutil"
	"k8s.io/apimachinery/pkg/api/errors"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (handler *APIHandler) ListApp(w http.ResponseWriter, r *http.Request) {
	fal := &fusionappv1alpha1.FusionAppList{}
	err := handler.client.List(context.TODO(), &client.ListOptions{}, fal)
	if err != nil {
		log.Warningf("failed to list apps", err)
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
		return
	}
	apps := make([]types.App, 0)
	for _, fusionApp := range fal.Items {
		apps = append(apps, *types.V1alpha1AppToApp(&fusionApp))
	}
	responseJSON(apps, w, http.StatusOK)
}

func (handler *APIHandler) CreateApp(w http.ResponseWriter, r *http.Request) {
	appAPICreateBody := new(types.AppAPICreateBody)
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	defer r.Body.Close()
	if err != nil {
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
		return
	}
	if err := json.Unmarshal(body, &appAPICreateBody); err != nil {
		if nerr := json.NewEncoder(w).Encode(err); nerr != nil {
			responseJSON(Message{nerr.Error()}, w, http.StatusUnprocessableEntity)
		} else {
			responseJSON(Message{err.Error()}, w, http.StatusBadRequest)
		}
		return
	}
	app := appAPICreateBody.AppSpec
	fusionApp := types.AppToV1alpha1App(&app)
	err = handler.client.Create(context.TODO(), fusionApp)
	if err != nil {
		log.Warningf("Failed to create App %v: %v", app.Name, err)
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	} else {
		responseJSON("ok", w, http.StatusOK)
	}
}

func (handler *APIHandler) UpdateApp(w http.ResponseWriter, r *http.Request) {
	appAPIPutBody := new(types.AppAPIPutBody)
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	defer r.Body.Close()
	if err != nil {
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
		return
	}

	if err := json.Unmarshal(body, &appAPIPutBody); err != nil {
		if nerr := json.NewEncoder(w).Encode(err); nerr != nil {
			responseJSON(Message{nerr.Error()}, w, http.StatusUnprocessableEntity)
		} else {
			responseJSON(Message{err.Error()}, w, http.StatusBadRequest)
		}
		return
	}
	name := appAPIPutBody.Name
	fusionApp := new(fusionappv1alpha1.FusionApp)
	err = handler.client.Get(context.TODO(), client.ObjectKey{
		Name: name}, fusionApp)
	if errors.IsNotFound(err) {
		err := fmt.Errorf("app \"%s\" not exists", name)
		responseJSON(Message{err.Error()}, w, http.StatusNotFound)
		return
	} else if err != nil {
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
		return
	}
	types.UpdateAppWithAppSpec(fusionApp, &appAPIPutBody.AppSpec)
	err = handler.client.Update(context.TODO(), fusionApp)
	if err != nil {
		log.Warningf("Failed to update fusionapp %v: %v", fusionApp.Name, err)
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	} else {
		responseJSON("ok", w, http.StatusOK)
	}
}

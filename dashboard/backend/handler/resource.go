package handler

import (
	"context"
	resourcev1alpha1 "github.com/fusion-app/fusion-app/pkg/apis/fusionapp/v1alpha1"
	log "github.com/sirupsen/logrus"
	"net/http"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

func (handler *APIHandler) ListResources(w http.ResponseWriter, r *http.Request) {
	dsl := &resourcev1alpha1.ResourceList{}

	err := handler.client.List(context.TODO(), &client.ListOptions{}, dsl)

	if err != nil {
		log.Warningf("failed to list datasets: %v", err)
		responseJSON(Message{err.Error()}, w, http.StatusInternalServerError)
	} else {
		responseJSON(ResourceList{Resources: dsl.Items}, w, http.StatusOK)
	}
}

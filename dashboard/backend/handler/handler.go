package handler

import (
	"encoding/json"
	"fmt"
	"github.com/fusion-app/fusion-app/dashboard/backend/types"
	"github.com/fusion-app/fusion-app/pkg/apis"
	"github.com/fusion-app/fusion-app/pkg/util/k8sutil"
	"github.com/sirupsen/logrus"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"net/http"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/client/config"
)

type APIHandler struct {
	frontDir string

	resourcesNamespace string

	kubeConfig    *rest.Config
	client        client.Client
	kubeClient    kubernetes.Interface

}

func NewAPIHandler(frontDir string) (*APIHandler, error) {
	kubeConfig, err := config.GetConfig()
	if err != nil {
		return nil, err
	}
	logrus.Infof("Using static directory %s", frontDir)

	// setup client set
	clientset, err := setupClient(kubeConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to setup kubernetes client: %v", err)
	}

	// setup kubernetes rest client
	kubeClient, err := k8sutil.NewKubeClient()
	if err != nil {
		return nil, fmt.Errorf("failed to setup kubernetes client: %v", err)
	}

	apiHandler := &APIHandler{
		frontDir:      frontDir,
		client:        clientset,
		kubeClient:    kubeClient,
		kubeConfig:    kubeConfig,
	}

	// Set resources namespace
	apiHandler.resourcesNamespace = os.Getenv("RESOURCES_NAMESPACE")
	if len(apiHandler.resourcesNamespace) == 0 {
		apiHandler.resourcesNamespace = types.DefaultNamespace
	}

	return apiHandler, nil
}

func setupClient(config *rest.Config) (client.Client, error) {
	scheme := runtime.NewScheme()
	for _, addToSchemeFunc := range []func(s *runtime.Scheme) error{
		apis.AddToScheme,
		v1.AddToScheme,
	} {
		if err := addToSchemeFunc(scheme); err != nil {
			return nil, err
		}
	}

	clientset, err := client.New(config, client.Options{Scheme: scheme})
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

type Message struct {
	Message string `json:"message"`
}

func responseJSON(body interface{}, w http.ResponseWriter, statusCode int) {
	jsonResponse, err := json.Marshal(body)

	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_, _ = w.Write(jsonResponse)
}

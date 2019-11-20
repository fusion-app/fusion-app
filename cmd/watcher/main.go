package main

import (
	"encoding/json"
	"flag"
	"github.com/fusion-app/fusion-app/dashboard/backend/types"
	"github.com/fusion-app/fusion-app/pkg/apis/fusionapp/v1alpha1"
	fusionappclient "github.com/fusion-app/fusion-app/pkg/client/clientset/versioned"
	"github.com/fusion-app/fusion-app/pkg/util/k8sutil"
	"github.com/jcuga/golongpoll"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
	"net/http"
)

func main() {
	var ns string
	flag.StringVar(&ns, "namespace", "", "namespace")
	flag.Parse()

	// bootstrap config
	cfg, err := k8sutil.GetClusterConfig()
	if err != nil {
		panic(err)
	}
	// create the clientset
	clientset, err := fusionappclient.NewForConfig(cfg)
	if err != nil {
		log.Fatal(err)
	}
	manager, err := golongpoll.StartLongpoll(golongpoll.Options{
		LoggingEnabled: true,
		// NOTE: if not defined here, other options have reasonable defaults,
		// so no need specifying options you don't care about
	})
	if err != nil {
		log.Fatalf("Failed to create manager: %q", err)
	}
	go watchResourcesHandler(manager, clientset, ns)
	http.HandleFunc("/events", manager.SubscriptionHandler)
	_ = http.ListenAndServe("127.0.0.1:8081", nil)
}

// printPVCs prints a list of PersistentVolumeClaim on console
func watchResourcesHandler(manager *golongpoll.LongpollManager, clientset *fusionappclient.Clientset, ns string) {
	// watch future changes to Resources
	watcher, err := clientset.FusionappV1alpha1().Resources(ns).Watch(metav1.ListOptions{})
	if err != nil {
		log.Fatal(err)
	}
	ch := watcher.ResultChan()

	for event := range ch {
		resource, ok := event.Object.(*v1alpha1.Resource)
		if !ok {
			log.Fatal("unexpected type")
			continue
		}
		switch event.Type {
		case watch.Error:
			log.Printf("watcher error encountered\n", resource.GetName())
		default:
			data, _ := json.Marshal(&Message{Type: event.Type, Resource: *types.V1alpha1ResourceToResource(resource)})
			_ = manager.Publish("resources", string(data))
		}
	}
	watcher.Stop()
}
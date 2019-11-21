package main

import (
	"flag"
	fusionappclient "github.com/fusion-app/fusion-app/pkg/client/clientset/versioned"
	"github.com/fusion-app/fusion-app/pkg/util/k8sutil"
	"github.com/jcuga/golongpoll"
	log "github.com/sirupsen/logrus"

	"net/http"
	"strconv"
)

const (
	DefaultListenPort int    = 8081
)

var (
	listenPort  int
)

func main() {
	var ns string
	flag.StringVar(&ns, "namespace", "", "namespace")
	flag.IntVar(&listenPort, "port", DefaultListenPort, `port this server listen to`)
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
	p := ":" + strconv.Itoa(listenPort)
	http.HandleFunc("/events", manager.SubscriptionHandler)
	go watchResourcesHandler(manager, clientset, ns)
	go watchAppInstanceHandler(manager, clientset, ns)
	_ = http.ListenAndServe(p, nil)
}

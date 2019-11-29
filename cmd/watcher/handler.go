package main

import (
	"encoding/json"
	"github.com/fusion-app/fusion-app/dashboard/backend/types"
	"github.com/fusion-app/fusion-app/pkg/apis/fusionapp/v1alpha1"
	fusionappclient "github.com/fusion-app/fusion-app/pkg/client/clientset/versioned"
	"github.com/jcuga/golongpoll"
	log "github.com/sirupsen/logrus"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/watch"
)

func watchResourcesHandler(manager *golongpoll.LongpollManager, clientset *fusionappclient.Clientset, ns string) {
	for {
		rsl, err := clientset.FusionappV1alpha1().Resources(ns).List(metav1.ListOptions{})
		if err != nil {
			log.Print(err)
			continue
		}
		originalCount := len(rsl.Items)
		// watch future changes to Resources
		watcher, err := clientset.FusionappV1alpha1().Resources(ns).Watch(metav1.ListOptions{})
		if err != nil {
			log.Print(err)
			continue
		}
		for {
			handled := true
			select {
			case event := <-watcher.ResultChan():
				resource, ok := event.Object.(*v1alpha1.Resource)
				if !ok {
					log.Printf("unexpected type, expecting Resource")
					handled = false
				} else {
					switch event.Type {
					case watch.Error:
						log.Printf("watcher error encountered\n", resource.GetName())
						handled = false
					default:
						rs, modified := types.V1alpha1ResourceToResource(resource)
						if len(rs.Phase) == 0 {
							rs.Phase = v1alpha1.ProbePhaseNotReady
						}
						data, _ := json.Marshal(&ResourceMessage{Type: event.Type, Resource: *rs})
						if originalCount <= 0 {
							_ = manager.Publish("resources", string(data))
						} else {
							originalCount --
						}
						if modified {
							go clientset.FusionappV1alpha1().Resources(ns).Update(resource.DeepCopy())
						}
					}
				}
			default:
				break
			}
			if !handled {
				watcher.Stop()
				break
			}
		}
	}
}

func watchAppInstanceHandler(manager *golongpoll.LongpollManager, clientset *fusionappclient.Clientset, ns string) {
	for {
		fusionAppInstanceList, err := clientset.FusionappV1alpha1().FusionAppInstances(ns).List(metav1.ListOptions{})
		if err != nil {
			log.Print(err)
			continue
		}
		originalCount := len(fusionAppInstanceList.Items)
		// watch future changes to AppInstances
		watcher, err := clientset.FusionappV1alpha1().FusionAppInstances(ns).Watch(metav1.ListOptions{})
		if err != nil {
			log.Print(err)
			continue
		}
		for {
			handled := true
			select {
			case event := <-watcher.ResultChan():
				appInstance, ok := event.Object.(*v1alpha1.FusionAppInstance)
				if !ok {
					log.Printf("unexpected type, expecting FusionAppInstance")
					handled = false
				} else {
					switch event.Type {
					case watch.Error:
						log.Printf("watcher error encountered\n", appInstance.GetName())
						handled = false
					default:
						data, _ := json.Marshal(&FusionAppInstanceMessage{Type: event.Type, AppInstance: *types.V1alpha1AppInstanceToAppInstance(appInstance)})
						if originalCount <= 0 {
							_ = manager.Publish("appinstances", string(data))
						} else {
							originalCount --
						}
					}
				}
			default:
				break
			}
			if !handled {
				watcher.Stop()
				break
			}
		}
	}
}
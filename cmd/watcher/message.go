package main

import (
	"github.com/fusion-app/fusion-app/dashboard/backend/types"
	"k8s.io/apimachinery/pkg/watch"
)

type Message struct {
	Type       watch.EventType    `json:"type"`
	Resource   types.Resource     `json:"resource"`
}

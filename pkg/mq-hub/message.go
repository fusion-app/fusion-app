package mqhub

import "time"

type Message struct {
	Target TargetCRDOption `json:"Target"`
	UpdatePatch []byte `json:"UpdatePatch"`
	ProbeTime time.Time `json:"ProbeTime"`
}

type TargetCRDOption struct {
	UID string `json:"UID,omitempty"`
	Kind string `json:"Kind"`
	Name string `json:"Name"`
	Namesapce string `json:"Namespace"`
}
package mqhub

import "time"

type Message struct {
	Target      TargetCRDOption `json:"target"`
	UpdatePatch []PatchItem     `json:"updatePatch"`
	ProbeTime   time.Time       `json:"probeTime"`
}

type TargetCRDOption struct {
	UID       string `json:"uid"`
	Kind      string `json:"kind"`
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
}

type PatchItem struct {
	Op    PatchOperation `json:"op"`
	Path  string         `json:"path"`
	From  string         `json:"from,omitempty"`
	Value interface{}    `json:"value,omitempty"`
}

type PatchOperation string

const (
	Add     PatchOperation = "add"
	Remove  PatchOperation = "remove"
	Replace PatchOperation = "replace"
	Copy    PatchOperation = "copy"
	Move    PatchOperation = "move"
	Test    PatchOperation = "test"
)

package controller

import (
	"github.com/fusion-app/fusion-app/pkg/controller/resourceclaim"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, resourceclaim.Add)
}

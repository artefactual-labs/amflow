package controllers

import (
	"github.com/goadesign/goa"
)

// WebController implements the public resource.
type WebController struct {
	*goa.Controller
}

// NewWeb creates a public controller.
func NewWeb(service *goa.Service) *WebController {
	return &WebController{Controller: service.NewController("WebController")}
}

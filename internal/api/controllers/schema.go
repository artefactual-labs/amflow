package controllers

import (
	"github.com/goadesign/goa"
)

// SchemaController implements the schema resource.
type SchemaController struct {
	*goa.Controller
}

// NewSchemaController creates a schema controller.
func NewSchemaController(service *goa.Service) *SchemaController {
	return &SchemaController{Controller: service.NewController("SchemaController")}
}

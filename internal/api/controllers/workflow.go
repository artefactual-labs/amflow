package controllers

import (
	"time"

	"github.com/goadesign/goa"

	"github.com/sevein/amflow/internal/api/app"
	"github.com/sevein/amflow/internal/graph"
)

// WorkflowController implements the workflow resource.
type WorkflowController struct {
	*goa.Controller

	graph *graph.Workflow
}

// NewWorkflow creates a workflow controller.
func NewWorkflow(service *goa.Service, graph *graph.Workflow) *WorkflowController {
	return &WorkflowController{
		Controller: service.NewController("AmflowController"),
		graph:      graph,
	}
}

func (c *WorkflowController) AddLink(ctx *app.AddLinkWorkflowContext) error {
	return nil
}

func (c *WorkflowController) DeleteLink(ctx *app.DeleteLinkWorkflowContext) error {
	return nil
}

func (c *WorkflowController) MoveLink(ctx *app.MoveLinkWorkflowContext) error {
	return nil
}

var blob []byte

func (c *WorkflowController) Show(ctx *app.ShowWorkflowContext) (err error) {
	start := time.Now()
	// Avoid generation if it's already cached.
	if blob == nil {
		blob, err = c.graph.SVG()
	}
	elapsed := time.Since(start)
	c.Service.LogInfo("svg", "elapsed", elapsed)
	if err != nil {
		ctx.ResponseWriter.WriteHeader(500)
		return err
	}
	ctx.ResponseWriter.Write(blob)
	return nil
}

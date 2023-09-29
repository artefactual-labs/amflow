// Code generated by goagen v1.4.3, DO NOT EDIT.
//
// API "amflow": Application Contexts
//
// Command:
// $ goagen
// --design=github.com/artefactual-labs/amflow/design
// --out=/home/jesus/Projects/amflow/internal/api
// --version=v1.4.3

package app

import (
	"context"
	"github.com/goadesign/goa"
	"net/http"
)

// AddLinkWorkflowContext provides the workflow addLink action context.
type AddLinkWorkflowContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
	WorkflowID string
}

// NewAddLinkWorkflowContext parses the incoming request URL and body, performs validations and creates the
// context used by the workflow controller addLink action.
func NewAddLinkWorkflowContext(ctx context.Context, r *http.Request, service *goa.Service) (*AddLinkWorkflowContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	req.Request = r
	rctx := AddLinkWorkflowContext{Context: ctx, ResponseData: resp, RequestData: req}
	paramWorkflowID := req.Params["workflowID"]
	if len(paramWorkflowID) > 0 {
		rawWorkflowID := paramWorkflowID[0]
		rctx.WorkflowID = rawWorkflowID
	}
	return &rctx, err
}

// OK sends a HTTP response with status code 200.
func (ctx *AddLinkWorkflowContext) OK(resp []byte) error {
	if ctx.ResponseData.Header().Get("Content-Type") == "" {
		ctx.ResponseData.Header().Set("Content-Type", "text/plain")
	}
	ctx.ResponseData.WriteHeader(200)
	_, err := ctx.ResponseData.Write(resp)
	return err
}

// NotFound sends a HTTP response with status code 404.
func (ctx *AddLinkWorkflowContext) NotFound() error {
	ctx.ResponseData.WriteHeader(404)
	return nil
}

// DeleteLinkWorkflowContext provides the workflow deleteLink action context.
type DeleteLinkWorkflowContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
	LinkID     string
	WorkflowID string
}

// NewDeleteLinkWorkflowContext parses the incoming request URL and body, performs validations and creates the
// context used by the workflow controller deleteLink action.
func NewDeleteLinkWorkflowContext(ctx context.Context, r *http.Request, service *goa.Service) (*DeleteLinkWorkflowContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	req.Request = r
	rctx := DeleteLinkWorkflowContext{Context: ctx, ResponseData: resp, RequestData: req}
	paramLinkID := req.Params["linkID"]
	if len(paramLinkID) > 0 {
		rawLinkID := paramLinkID[0]
		rctx.LinkID = rawLinkID
	}
	paramWorkflowID := req.Params["workflowID"]
	if len(paramWorkflowID) > 0 {
		rawWorkflowID := paramWorkflowID[0]
		rctx.WorkflowID = rawWorkflowID
	}
	return &rctx, err
}

// OK sends a HTTP response with status code 200.
func (ctx *DeleteLinkWorkflowContext) OK(resp []byte) error {
	if ctx.ResponseData.Header().Get("Content-Type") == "" {
		ctx.ResponseData.Header().Set("Content-Type", "text/plain")
	}
	ctx.ResponseData.WriteHeader(200)
	_, err := ctx.ResponseData.Write(resp)
	return err
}

// NotFound sends a HTTP response with status code 404.
func (ctx *DeleteLinkWorkflowContext) NotFound() error {
	ctx.ResponseData.WriteHeader(404)
	return nil
}

// MoveLinkWorkflowContext provides the workflow moveLink action context.
type MoveLinkWorkflowContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
	LinkID     string
	WorkflowID string
}

// NewMoveLinkWorkflowContext parses the incoming request URL and body, performs validations and creates the
// context used by the workflow controller moveLink action.
func NewMoveLinkWorkflowContext(ctx context.Context, r *http.Request, service *goa.Service) (*MoveLinkWorkflowContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	req.Request = r
	rctx := MoveLinkWorkflowContext{Context: ctx, ResponseData: resp, RequestData: req}
	paramLinkID := req.Params["linkID"]
	if len(paramLinkID) > 0 {
		rawLinkID := paramLinkID[0]
		rctx.LinkID = rawLinkID
	}
	paramWorkflowID := req.Params["workflowID"]
	if len(paramWorkflowID) > 0 {
		rawWorkflowID := paramWorkflowID[0]
		rctx.WorkflowID = rawWorkflowID
	}
	return &rctx, err
}

// OK sends a HTTP response with status code 200.
func (ctx *MoveLinkWorkflowContext) OK(resp []byte) error {
	if ctx.ResponseData.Header().Get("Content-Type") == "" {
		ctx.ResponseData.Header().Set("Content-Type", "text/plain")
	}
	ctx.ResponseData.WriteHeader(200)
	_, err := ctx.ResponseData.Write(resp)
	return err
}

// NotFound sends a HTTP response with status code 404.
func (ctx *MoveLinkWorkflowContext) NotFound() error {
	ctx.ResponseData.WriteHeader(404)
	return nil
}

// ShowWorkflowContext provides the workflow show action context.
type ShowWorkflowContext struct {
	context.Context
	*goa.ResponseData
	*goa.RequestData
	WorkflowID string
}

// NewShowWorkflowContext parses the incoming request URL and body, performs validations and creates the
// context used by the workflow controller show action.
func NewShowWorkflowContext(ctx context.Context, r *http.Request, service *goa.Service) (*ShowWorkflowContext, error) {
	var err error
	resp := goa.ContextResponse(ctx)
	resp.Service = service
	req := goa.ContextRequest(ctx)
	req.Request = r
	rctx := ShowWorkflowContext{Context: ctx, ResponseData: resp, RequestData: req}
	paramWorkflowID := req.Params["workflowID"]
	if len(paramWorkflowID) > 0 {
		rawWorkflowID := paramWorkflowID[0]
		rctx.WorkflowID = rawWorkflowID
	}
	return &rctx, err
}

// OK sends a HTTP response with status code 200.
func (ctx *ShowWorkflowContext) OK(resp []byte) error {
	if ctx.ResponseData.Header().Get("Content-Type") == "" {
		ctx.ResponseData.Header().Set("Content-Type", "application/xml")
	}
	ctx.ResponseData.WriteHeader(200)
	_, err := ctx.ResponseData.Write(resp)
	return err
}

// NotFound sends a HTTP response with status code 404.
func (ctx *ShowWorkflowContext) NotFound() error {
	ctx.ResponseData.WriteHeader(404)
	return nil
}

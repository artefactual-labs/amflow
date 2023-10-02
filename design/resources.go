package design

import (
	. "github.com/goadesign/goa/design"
	. "github.com/goadesign/goa/design/apidsl"
)

var _ = Resource("workflow", func() {
	BasePath("/workflow")
	Origin("*", func() {
		Methods("GET", "POST", "PUT", "PATCH", "DELETE")
		MaxAge(600)
	})

	Action("show", func() {
		Description("Read workflow")
		Routing(GET("/:workflowID"))
		Params(func() {
			Param("workflowID", String, "Workflow ID")
		})
		Response(OK, "application/xml")
		Response(NotFound)
	})

	Action("addLink", func() {
		Description("Add link")
		Routing(PATCH("/:workflowID/links"))
		Params(func() {
			Param("workflowID", String, "Workflow ID")
		})
		Response(OK)
		Response(NotFound)
	})

	Action("moveLink", func() {
		Description("Move link")
		Routing(PATCH("/:workflowID/links/:linkID"))
		Params(func() {
			Param("workflowID", String, "Workflow ID")
			Param("linkID", String, "Link ID")
		})
		Response(OK)
		Response(NotFound)
	})

	Action("deleteLink", func() {
		Description("Delete link")
		Routing(DELETE("/:workflowID/links/:linkID"))
		Params(func() {
			Param("workflowID", String, "Workflow ID")
			Param("linkID", String, "Link ID")
		})
		Response(OK)
		Response(NotFound)
	})
})

var _ = Resource("web", func() {
	Origin("*", func() {
		Methods("GET, OPTIONS")
	})
	Files("/*filepath", "dist")
})

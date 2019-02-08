package api

import (
	"net/http"
	"time"

	"github.com/goadesign/goa"
	logadapter "github.com/goadesign/goa/logging/logrus"
	"github.com/goadesign/goa/middleware"
	"github.com/gobuffalo/packr/v2"
	"github.com/sirupsen/logrus"

	"github.com/sevein/amflow/internal/api/app"
	"github.com/sevein/amflow/internal/api/controllers"
	"github.com/sevein/amflow/internal/graph"
)

func Create(graph *graph.Workflow, logger *logrus.Entry) *goa.Service {
	service := goa.New("amflow")
	service.Use(middleware.RequestID())
	service.Use(middleware.ErrorHandler(service, true))
	service.Use(middleware.Recover())
	service.WithLogger(logadapter.FromEntry(logger))

	// Sane server timeouts.
	service.Server.ReadTimeout = 5 * time.Second
	service.Server.WriteTimeout = 10 * time.Second

	// Workflow controller.
	app.MountWorkflowController(service, controllers.NewWorkflow(service, graph))

	// Schema controller.
	schemaCtrl := controllers.NewSchemaController(service)
	schemaCtrl.FileSystem = func(dir string) http.FileSystem {
		return packr.New("schema", "../../public/schema")
	}
	app.MountSchemaController(service, schemaCtrl)

	// Swagger controller.
	swaggerCtrl := controllers.NewSwaggerController(service)
	swaggerCtrl.FileSystem = func(dir string) http.FileSystem {
		return packr.New("swagger", "../../public/swagger")
	}
	app.MountSwaggerController(service, swaggerCtrl)

	// Web controller.
	webCtrl := controllers.NewSwaggerController(service)
	webCtrl.FileSystem = func(dir string) http.FileSystem {
		return packr.New("assets", "../../public/web")
	}
	app.MountWebController(service, webCtrl)

	return service
}

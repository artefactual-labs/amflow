package api

import (
	"io/fs"
	"net/http"
	"time"

	"github.com/goadesign/goa"
	logadapter "github.com/goadesign/goa/logging/logrus"
	"github.com/goadesign/goa/middleware"
	"github.com/sirupsen/logrus"

	"github.com/artefactual-labs/amflow/internal/api/app"
	"github.com/artefactual-labs/amflow/internal/api/controllers"
	"github.com/artefactual-labs/amflow/internal/graph"
	"github.com/artefactual-labs/amflow/public"
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
	wfCtrl := controllers.NewWorkflow(service, graph)
	app.MountWorkflowController(service, wfCtrl)

	// Web controller.
	webCtrl := controllers.NewSwaggerController(service)
	webCtrl.FileSystem = func(dir string) http.FileSystem {
		assetsDir, _ := fs.Sub(public.Assets, "web")
		return http.FS(assetsDir)
	}
	app.MountWebController(service, webCtrl)

	return service
}

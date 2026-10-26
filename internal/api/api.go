package api

import (
	"bytes"
	"io/fs"
	"net/http"
	"sync"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/artefactual-labs/amflow/internal/graph"
	"github.com/artefactual-labs/amflow/web"
)

// Create returns an HTTP server for the workflow and embedded web application.
func Create(workflow *graph.Workflow, logger *logrus.Entry) *http.Server {
	assets, err := fs.Sub(web.Assets, "dist")
	if err != nil {
		panic(err)
	}

	var (
		mu   sync.Mutex
		blob []byte
	)
	loadSVG := func() ([]byte, error) {
		mu.Lock()
		defer mu.Unlock()
		if blob != nil {
			return blob, nil
		}
		svg, err := workflow.SVG(&graph.VizOpts{})
		if err != nil {
			return nil, err
		}
		blob = svg
		return blob, nil
	}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /workflow/{workflowID}", func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		svg, err := loadSVG()
		logger.WithField("elapsed", time.Since(start)).Info("svg")
		if err != nil {
			logger.WithError(err).Error("Rendering workflow")
			http.Error(w, "Could not render workflow", http.StatusInternalServerError)
			return
		}
		w.Header().Set("Content-Type", "image/svg+xml")
		http.ServeContent(w, r, "workflow.svg", time.Time{}, bytes.NewReader(svg))
	})

	// These routes retain the existing placeholders for link editing.
	empty := func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}
	mux.HandleFunc("PATCH /workflow/{workflowID}/links", empty)
	mux.HandleFunc("PATCH /workflow/{workflowID}/links/{linkID}", empty)
	mux.HandleFunc("DELETE /workflow/{workflowID}/links/{linkID}", empty)

	workflowOptions := func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, DELETE")
		w.Header().Set("Access-Control-Max-Age", "600")
		w.WriteHeader(http.StatusOK)
	}
	mux.HandleFunc("OPTIONS /workflow/{workflowID}", workflowOptions)
	mux.HandleFunc("OPTIONS /workflow/{workflowID}/links", workflowOptions)
	mux.HandleFunc("OPTIONS /workflow/{workflowID}/links/{linkID}", workflowOptions)
	mux.HandleFunc("OPTIONS /", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.WriteHeader(http.StatusOK)
	})
	mux.Handle("GET /", http.FileServerFS(assets))

	return &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if recovered := recover(); recovered != nil {
					if recovered == http.ErrAbortHandler {
						panic(recovered)
					}
					logger.WithField("panic", recovered).Error("Handling HTTP request")
					http.Error(w, "Internal server error", http.StatusInternalServerError)
				}
			}()
			if origin := r.Header.Get("Origin"); origin != "" {
				w.Header().Set("Access-Control-Allow-Origin", origin)
				w.Header().Add("Vary", "Origin")
			}
			mux.ServeHTTP(w, r)
		}),
		ReadTimeout:       5 * time.Second,
		ReadHeaderTimeout: 5 * time.Second,
		WriteTimeout:      10 * time.Second,
	}
}

package api_test

import (
	"encoding/xml"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"

	"github.com/sirupsen/logrus"
	"gotest.tools/v3/assert"
	"gotest.tools/v3/assert/cmp"

	"github.com/artefactual-labs/amflow/internal/api"
	"github.com/artefactual-labs/amflow/internal/graph"
	amjson "github.com/artefactual-labs/amflow/internal/graph/encoding"
	"github.com/artefactual-labs/amflow/web"
)

func newServer(t *testing.T, description string) *http.Server {
	t.Helper()
	data, err := amjson.LoadWorkflowData([]byte(fmt.Sprintf(`{
		"links": {"test-link": {"description": {"en": %q}, "start": true}}
	}`, description)))
	assert.NilError(t, err)
	logger := logrus.New()
	logger.SetOutput(io.Discard)
	return api.Create(graph.New(data), logrus.NewEntry(logger))
}

func request(handler http.Handler, method, path string) *httptest.ResponseRecorder {
	recorder := httptest.NewRecorder()
	handler.ServeHTTP(recorder, httptest.NewRequest(method, path, nil))
	return recorder
}

func TestWebAssets(t *testing.T) {
	server := newServer(t, "Workflow")
	index, err := fs.ReadFile(web.Assets, "dist/index.html")
	assert.NilError(t, err)
	response := request(server.Handler, http.MethodGet, "/")
	assert.Equal(t, response.Code, http.StatusOK)
	assert.Equal(t, response.Header().Get("Content-Type"), "text/html; charset=utf-8")
	assert.DeepEqual(t, response.Body.Bytes(), index)

	err = fs.WalkDir(web.Assets, "dist/assets", func(path string, entry fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if entry.IsDir() {
			return nil
		}
		want, err := fs.ReadFile(web.Assets, path)
		assert.NilError(t, err)
		response := request(server.Handler, http.MethodGet, strings.TrimPrefix(path, "dist"))
		assert.Equal(t, response.Code, http.StatusOK, path)
		assert.DeepEqual(t, response.Body.Bytes(), want)
		return nil
	})
	assert.NilError(t, err)

	response = request(server.Handler, http.MethodGet, "/missing.js")
	assert.Equal(t, response.Code, http.StatusNotFound)
	response = request(server.Handler, http.MethodPost, "/")
	assert.Equal(t, response.Code, http.StatusMethodNotAllowed)
}

func TestWorkflowSVGConcurrentRequests(t *testing.T) {
	if testing.Short() {
		t.Skip("requires Graphviz")
	}
	server := newServer(t, "Concurrent workflow")
	responses := make([]*httptest.ResponseRecorder, 8)
	var wg sync.WaitGroup
	start := make(chan struct{})
	for i := range responses {
		wg.Go(func() {
			<-start
			responses[i] = request(server.Handler, http.MethodGet, "/workflow/default")
		})
	}
	close(start)
	wg.Wait()

	for _, response := range responses {
		assert.Equal(t, response.Code, http.StatusOK)
		assert.Equal(t, response.Header().Get("Content-Type"), "image/svg+xml")
		assert.DeepEqual(t, response.Body.Bytes(), responses[0].Body.Bytes())
	}
	var svg struct {
		XMLName xml.Name `xml:"svg"`
	}
	assert.NilError(t, xml.Unmarshal(responses[0].Body.Bytes(), &svg))
	assert.Assert(t, cmp.Contains(responses[0].Body.String(), "Concurrent workflow"))

	response := request(server.Handler, http.MethodHead, "/workflow/default")
	assert.Equal(t, response.Code, http.StatusOK)
	assert.Equal(t, response.Body.Len(), 0)
	assert.Equal(t, response.Header().Get("Content-Length"), responses[0].Header().Get("Content-Length"))
}

func TestWorkflowSVGCacheIsScopedToServer(t *testing.T) {
	if testing.Short() {
		t.Skip("requires Graphviz")
	}
	first := newServer(t, "First workflow")
	second := newServer(t, "Second workflow")
	for _, tc := range []struct {
		server      *http.Server
		description string
	}{
		{first, "First workflow"},
		{second, "Second workflow"},
		{first, "First workflow"},
	} {
		response := request(tc.server.Handler, http.MethodGet, "/workflow/default")
		assert.Equal(t, response.Code, http.StatusOK)
		assert.Assert(t, cmp.Contains(response.Body.String(), tc.description))
	}
}

func TestWorkflowSVGRenderFailureAndRetry(t *testing.T) {
	if testing.Short() {
		t.Skip("requires Graphviz")
	}
	server := newServer(t, "Retry workflow")
	path := os.Getenv("PATH")
	t.Setenv("PATH", t.TempDir())
	response := request(server.Handler, http.MethodGet, "/workflow/default")
	assert.Equal(t, response.Code, http.StatusInternalServerError)
	assert.Equal(t, response.Body.String(), "Could not render workflow\n")

	// A failed render must not poison subsequent requests.
	t.Setenv("PATH", path)
	response = request(server.Handler, http.MethodGet, "/workflow/default")
	assert.Equal(t, response.Code, http.StatusOK)
	assert.Assert(t, cmp.Contains(response.Body.String(), "Retry workflow"))

	// A successful render remains available without running Graphviz again.
	t.Setenv("PATH", t.TempDir())
	cached := request(server.Handler, http.MethodGet, "/workflow/default")
	assert.Equal(t, cached.Code, http.StatusOK)
	assert.DeepEqual(t, cached.Body.Bytes(), response.Body.Bytes())
}

func TestWorkflowRoutes(t *testing.T) {
	server := newServer(t, "Workflow")
	for _, tc := range []struct {
		method string
		path   string
		status int
	}{
		{http.MethodPatch, "/workflow/default/links", http.StatusOK},
		{http.MethodPatch, "/workflow/default/links/link", http.StatusOK},
		{http.MethodDelete, "/workflow/default/links/link", http.StatusOK},
		{http.MethodPost, "/workflow/default", http.StatusMethodNotAllowed},
		{http.MethodGet, "/workflow/default/links/missing", http.StatusNotFound},
	} {
		t.Run(tc.method+" "+tc.path, func(t *testing.T) {
			response := request(server.Handler, tc.method, tc.path)
			assert.Equal(t, response.Code, tc.status)
		})
	}

	for _, path := range []string{
		"/workflow/default",
		"/workflow/default/links",
		"/workflow/default/links/link",
		"/",
	} {
		t.Run("OPTIONS "+path, func(t *testing.T) {
			r := httptest.NewRequest(http.MethodOptions, path, nil)
			r.Header.Set("Origin", "http://localhost:3000")
			r.Header.Set("Access-Control-Request-Method", http.MethodPatch)
			response := httptest.NewRecorder()
			server.Handler.ServeHTTP(response, r)
			assert.Equal(t, response.Code, http.StatusOK)
			assert.Equal(t, response.Header().Get("Access-Control-Allow-Origin"), "http://localhost:3000")
			assert.Equal(t, response.Header().Get("Vary"), "Origin")
			if strings.HasPrefix(path, "/workflow/") {
				assert.Assert(t, cmp.Contains(response.Header().Get("Access-Control-Allow-Methods"), "PATCH"))
				assert.Equal(t, response.Header().Get("Access-Control-Max-Age"), "600")
			}
		})
	}
}

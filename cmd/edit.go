package cmd

import (
	"io"
	"io/ioutil"
	"net"
	"os/exec"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/sevein/amflow/internal/api"
	"github.com/sevein/amflow/internal/graph"
	"github.com/sevein/amflow/internal/graph/encoding"
)

var (
	file, addr string
)

func newCmdEdit(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit",
		Short: "Edit the workflow",
		RunE: func(cmd *cobra.Command, args []string) error {
			return edit(out)
		},
	}
	cmd.Flags().StringVarP(&file, "file", "f", "", "Path of JSON-encoded workflow document")
	cmd.Flags().StringVarP(&addr, "addr", "", ":2323", "HTTP service address (default: ':2323')")
	return cmd
}

func edit(out io.Writer) error {
	var (
		bytes []byte
		err   error
	)
	if file == "" {
		logrus.WithFields(logrus.Fields{"mode": "embedded"}).Debug("Loading workfow")
		bytes, err = graph.WorkflowSchemaBox.Find("example.json")
	} else {
		logrus.WithFields(logrus.Fields{"mode": "file", "source": file}).Debug("Loading workfow")
		bytes, err = ioutil.ReadFile(file)
	}
	if err != nil {
		return err
	}

	// Decode and populate workflow.
	wd, err := encoding.LoadWorkflowData(bytes)
	if err != nil {
		return err
	}
	g := graph.New(wd)
	logrus.WithFields(logrus.Fields{
		"bytes":    len(bytes),
		"vertices": g.Nodes().Len(),
	}).Debug("Workflow loaded")

	// Find dot executable.
	if _, err = exec.LookPath("dot"); err != nil {
		logrus.Warn("dot (Graphviz) is not installed")
	}

	// Start API server.
	ln, err := getListener(addr)
	if err != nil {
		return err
	}
	defer ln.Close()
	logger := logrus.WithField("subsystem", "api")
	svc := api.Create(g, logger)
	if err := svc.Serve(ln); err != nil {
		svc.LogError("startup", "err", err)
	}

	return nil
}

func getListener(addr string) (*net.TCPListener, error) {
	if addr == "" {
		addr = "localhost:0"
	}
	tcpa, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	ln, err := net.ListenTCP("tcp", tcpa)
	if err != nil {
		return nil, err
	}
	return ln, nil
}

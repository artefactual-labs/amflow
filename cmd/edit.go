package cmd

import (
	"io"
	"net"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/artefactual-labs/amflow/internal/api"
)

var (
	editFile   string
	editAddr   string
	editLatest bool
)

func newCmdEdit(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit",
		Short: "Edit the workflow",
		RunE: func(cmd *cobra.Command, args []string) error {
			return edit(out)
		},
	}
	cmd.Flags().StringVarP(&editFile, "file", "f", "", "Path or URL of the JSON-encoded workflow document")
	cmd.Flags().StringVarP(&editAddr, "addr", "", ":2323", "HTTP service address (default: ':2323')")
	cmd.Flags().BoolVarP(&editLatest, "latest", "", false, "Download the latest workflow available in QA")
	return cmd
}

const latestWorkflow = "https://raw.githubusercontent.com/artefactual/archivematica/qa/1.x/src/MCPServer/lib/assets/workflow.json"

func edit(out io.Writer) error {
	if editLatest {
		editFile = latestWorkflow
	}
	w, err := load(editFile)
	if err != nil {
		return err
	}

	checkDot()

	// Start API server.
	ln, err := getListener(editAddr)
	if err != nil {
		return err
	}
	defer ln.Close()
	logger := logrus.WithField("subsystem", "api")
	logger.WithField("port", ln.Addr()).Info("Staring API server")
	svc := api.Create(w, logger)
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

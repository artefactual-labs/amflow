package cmd

import (
	"io"
	"net"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/sevein/amflow/internal/api"
)

var (
	editFile string
	editAddr string
)

func newCmdEdit(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "edit",
		Short: "Edit the workflow",
		RunE: func(cmd *cobra.Command, args []string) error {
			return edit(out)
		},
	}
	cmd.Flags().StringVarP(&editFile, "file", "f", "", "Path of JSON-encoded workflow document")
	cmd.Flags().StringVarP(&editAddr, "addr", "", ":2323", "HTTP service address (default: ':2323')")
	return cmd
}

func edit(out io.Writer) error {
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

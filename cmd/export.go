package cmd

import (
	"fmt"
	"io"
	"io/ioutil"
	"os/exec"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/sevein/amflow/internal/graph"
	"github.com/sevein/amflow/internal/graph/encoding"
)

var exportFile string

func newCmdExport(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export the workflow in DOT format",
		Long: `
qwefqwe
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return export(out)
		},
	}
	cmd.Flags().StringVarP(&exportFile, "file", "f", "", "Path of JSON-encoded workflow document")
	return cmd
}

func export(out io.Writer) error {
	var (
		bytes []byte
		err   error
	)
	if exportFile == "" {
		logrus.WithFields(logrus.Fields{"mode": "embedded"}).Debug("Loading workfow")
		bytes, err = graph.WorkflowSchemaBox.Find("example.json")
	} else {
		logrus.WithFields(logrus.Fields{"mode": "file", "source": exportFile}).Debug("Loading workfow")
		bytes, err = ioutil.ReadFile(exportFile)
	}
	if err != nil {
		return err
	}

	// Decode and populate workflow.
	wd, err := encoding.LoadWorkflowData(bytes)
	if err != nil {
		return err
	}
	w := graph.New(wd)
	logrus.WithFields(logrus.Fields{
		"bytes":    len(bytes),
		"vertices": w.Nodes().Len(),
	}).Debug("Workflow loaded")

	// Find dot executable.
	if _, err = exec.LookPath("dot"); err != nil {
		logrus.Warn("dot (Graphviz) is not installed")
	}

	// Print it out.
	blob, err := w.DOT()
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(out, string(blob))
	return err
}

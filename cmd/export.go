package cmd

import (
	"fmt"
	"io"

	"github.com/artefactual-labs/amflow/internal/graph"
	"github.com/spf13/cobra"
)

var (
	exportFile   string
	exportFormat string
	exportFull   bool
	exportLatest bool
)

func newCmdExport(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "export",
		Short: "Export the workflow in DOT format",
		Long: `
Export the workflow in DOT format.

Usage example:

    $ amflow export --file=test.json --format=dot | dot -v -Tsvg > test.svg
`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return export(out)
		},
	}
	cmd.Flags().StringVarP(&exportFile, "file", "f", "", "Path of JSON-encoded workflow document")
	cmd.Flags().StringVarP(&exportFormat, "format", "", "dot", "Format of the export")
	cmd.Flags().BoolVarP(&exportLatest, "latest", "", false, "Download the latest workflow available in QA")
	cmd.Flags().BoolVarP(&exportFull, "full", "", false, "Include the full graph (slower)")
	return cmd
}

func export(out io.Writer) error {
	if exportLatest {
		exportFile = latestWorkflow
	}
	w, err := load(exportFile)
	if err != nil {
		return err
	}

	checkDot()

	var blob []byte
	switch exportFormat {
	case "dot":
		blob, err = w.DOT(&graph.VizOpts{Full: exportFull})
	case "svg":
		blob, err = w.SVG(&graph.VizOpts{Full: exportFull})
	}
	if err != nil {
		return err
	}

	_, err = fmt.Fprint(out, string(blob))

	return err
}

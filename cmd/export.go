package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"
)

var (
	exportFile   string
	exportFormat string
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
	return cmd
}

func export(out io.Writer) error {
	w, err := load(exportFile)
	if err != nil {
		return err
	}

	checkDot()

	// Print it out.
	blob, err := w.DOT()
	if err != nil {
		return err
	}
	_, err = fmt.Fprint(out, string(blob))
	return err
}

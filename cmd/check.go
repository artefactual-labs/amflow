package cmd

import (
	"io"

	"github.com/spf13/cobra"
)

var (
	checkFile   string
	checkLatest bool
)

func newCmdCheck(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "check",
		Short: "Verify workflow integrity",
		RunE: func(cmd *cobra.Command, args []string) error {
			return check(out)
		},
	}
	cmd.Flags().StringVarP(&checkFile, "file", "f", "", "Path or URL of the JSON-encoded workflow document")
	cmd.Flags().BoolVarP(&checkLatest, "latest", "", false, "Download the latest workflow available in QA")
	return cmd
}

func check(out io.Writer) error {
	if checkLatest {
		checkFile = latestWorkflow
	}

	_, err := load(checkFile)

	return err
}

package cmd

import (
	"fmt"
	"io"

	"github.com/spf13/cobra"

	"github.com/artefactual-labs/amflow/internal/version"
)

func newCmdVersion(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Print the version information",
		RunE: func(cmd *cobra.Command, args []string) error {
			fmt.Println(version.Version())
			return nil
		},
	}
	return cmd
}

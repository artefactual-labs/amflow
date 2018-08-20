package cmd

import (
	"io"
	"os"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/sevein/amflow/internal/constants"
	"github.com/sevein/amflow/internal/version"
)

var (
	v string
)

var rootCmd = &cobra.Command{
	Use:   "amflow",
	Short: "A tool that facilitates workflow editing for Archivematica.",
}

func Run() error {
	c := command(os.Stdout, os.Stderr)
	return c.Execute()
}

func command(out, err io.Writer) *cobra.Command {
	rootCmd.PersistentPreRunE = func(cmd *cobra.Command, args []string) error {
		if err := setUpLogs(err, v); err != nil {
			return err
		}
		rootCmd.SilenceUsage = true
		logrus.Infof("amflow %+s", version.Get())
		return nil
	}

	rootCmd.SilenceErrors = true
	rootCmd.AddCommand(newCmdVersion(out))
	rootCmd.AddCommand(newCmdEdit(out))
	rootCmd.AddCommand(newCmdExport(out))

	rootCmd.PersistentFlags().StringVarP(&v, "verbosity", "v", constants.DefaultLogLevel.String(), "Log level (debug, info, warn, error, fatal, panic")

	return rootCmd
}

func setUpLogs(out io.Writer, level string) error {
	logrus.SetOutput(out)
	lvl, err := logrus.ParseLevel(v)
	if err != nil {
		return errors.Wrap(err, "parsing log level")
	}
	logrus.SetLevel(lvl)
	return nil
}

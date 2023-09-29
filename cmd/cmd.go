/*
Package cmd implements command-line interfaces.
*/
package cmd

import (
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"go.uber.org/multierr"

	"github.com/artefactual-labs/amflow/internal/graph"
	"github.com/artefactual-labs/amflow/internal/graph/encoding"
	"github.com/artefactual-labs/amflow/internal/version"
)

var (
	v               string
	defaultLogLevel = logrus.InfoLevel
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
		logrus.Infof("amflow %+s", version.Version())
		return nil
	}

	rootCmd.SilenceErrors = true
	rootCmd.AddCommand(newCmdVersion(out))
	rootCmd.AddCommand(newCmdEdit(out))
	rootCmd.AddCommand(newCmdExport(out))
	rootCmd.AddCommand(newCmdSearch(out))
	rootCmd.AddCommand(newCmdCheck(out))

	rootCmd.PersistentFlags().StringVarP(&v, "verbosity", "v", defaultLogLevel.String(), "Log level (debug, info, warn, error, fatal, panic")

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

// load returns the workflow loaded from the contents of the file.
func load(file string) (*graph.Workflow, error) {
	var (
		bytes []byte
		err   error
	)

	// Load workflow bytes.
	if file == "" {
		logrus.WithFields(logrus.Fields{"mode": "embedded"}).Info("Loading workfow")
		bytes = graph.WorkflowSample
	} else if isURL(file) {
		logrus.WithFields(logrus.Fields{"mode": "file", "source": file}).Info("Downloading workfow")
		bytes, err = downloadRemote(file)
	} else {
		logrus.WithFields(logrus.Fields{"mode": "file", "source": file}).Info("Loading workfow")
		bytes, err = os.ReadFile(file)
	}
	if err != nil {
		return nil, errors.WithMessage(err, "Workflow could not be retrieved")
	}

	// Decode it.
	data, err := encoding.LoadWorkflowData(bytes)
	if err != nil {
		return nil, err
	}

	// Populate it.
	w := graph.New(data)
	logrus.WithFields(logrus.Fields{
		"bytes":    len(bytes),
		"vertices": w.Nodes().Len(),
	}).Debug("Workflow loaded")

	// Check for errors.
	for _, err := range multierr.Errors(w.Check()) {
		logrus.WithFields(logrus.Fields{
			"err": err,
		}).Warn("Unhealthy workflow warning")
	}

	return w, nil
}

func isURL(addr string) bool {
	u, err := url.Parse(addr)
	if err != nil {
		return false
	}
	if u.Scheme == "" {
		return false
	}
	return true
}

func downloadRemote(addr string) ([]byte, error) {
	resp, err := http.Get(addr)
	if err != nil {
		return nil, errors.WithMessage(err, "remote resource could not be retrieved")
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		return nil, fmt.Errorf("remote server returned and unexpected response with status code: %d", resp.StatusCode)
	}
	defer resp.Body.Close()
	bytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.WithMessage(err, "remote resource could not be loaded")
	}
	return bytes, nil
}

func checkDot() {
	if _, err := exec.LookPath("dot"); err != nil {
		logrus.Warn("dot (Graphviz) is not installed")
	}
}

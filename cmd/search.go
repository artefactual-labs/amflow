package cmd

import (
	"errors"
	"fmt"
	"io"
	"text/tabwriter"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/sevein/amflow/internal/graph"
	amjson "github.com/sevein/amflow/internal/graph/encoding"
)

var (
	searchFile string
	searchTo   string
)

func newCmdSearch(out io.Writer) *cobra.Command {
	cmd := &cobra.Command{
		Use: "search",
		RunE: func(cmd *cobra.Command, args []string) error {
			return search(out)
		},
	}
	cmd.Flags().StringVarP(&searchFile, "file", "f", "", "Path of JSON-encoded workflow document")
	cmd.Flags().StringVarP(&searchTo, "to", "", "", "")
	return cmd
}

func search(out io.Writer) error {
	w, err := load(searchFile)
	if err != nil {
		return err
	}

	checkDot()

	if searchTo == "" {
		return errors.New("the target is not defined, use argument --to=<uuid>")
	}

	ints := graph.ListAncestors(w, searchTo)
	if ints.Count() == 0 {
		return errors.New("no results found")
	}

	wr := tabwriter.NewWriter(out, 0, 0, 3, ' ', tabwriter.AlignRight|tabwriter.Debug)
	for id := range ints {
		vertice := w.VertexByID(id)
		var t, amid, name, details string
		amid = vertice.AMID()
		switch v := vertice.(type) {
		case *graph.VertexChainLink:
			t = "Chain"
			_ch := vertice.Source().(*amjson.Chain)
			name = _ch.Description["en"]
		case *graph.VertexLink:
			t = "Link"
			_ln := vertice.Source().(*amjson.Link)
			name = _ln.Description["en"]
			details = _ln.Config.Execute
		case *graph.VertexWatcheDir:
			t = "Watched directory"
			_wd := vertice.Source().(*amjson.WatchedDirectory)
			name = _wd.Path
		default:
			logrus.Warnf("I don't know about %+v", v)
		}
		fmt.Fprintln(wr, fmt.Sprintf("\t%d\t%s\t%s\t%s\t%s\t", id, t, amid, name, details))
	}
	wr.Flush()

	return nil
}

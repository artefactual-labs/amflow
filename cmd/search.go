package cmd

import (
	"errors"
	"io"

	"github.com/olekukonko/tablewriter"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/sevein/amflow/internal/graph"
	amjson "github.com/sevein/amflow/internal/graph/encoding"
)

var (
	searchFile     string
	searchTo       string
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

	if searchTo == "" {
		return errors.New("the target is not defined, use argument --to=<uuid>")
	}

	ints := graph.ListAncestors(w, searchTo)
	if ints.Count() == 0 {
		return errors.New("no results found")
	}

	table := tablewriter.NewWriter(out)
	table.SetHeader([]string{"ID", "Type", "Module", "Description"})
	for id := range ints {
		vertice := w.VertexByID(id)
		var t, amid, desc, module string
		amid = vertice.AMID()
		switch v := vertice.(type) {
		case *graph.VertexChainLink:
			t = "Chain"
			_ch := vertice.Source().(*amjson.Chain)
			desc = _ch.Description["en"]
		case *graph.VertexLink:
			t = "Link"
			_ln := vertice.Source().(*amjson.Link)
			desc = _ln.Description["en"]
			module = _ln.Config.Execute
		case *graph.VertexWatcheDir:
			t = "Watched directory"
			_wd := vertice.Source().(*amjson.WatchedDirectory)
			desc = _wd.Path
		default:
			logrus.Warnf("I don't know about %+v", v)
		}
		table.Append([]string{
			truncate(amid, 36),
			t,
			truncate(module, 30),
			truncate(desc, 30),
		})
	}
	table.Render()

	return nil
}

func truncate(input string, size int) string {
	const dots = "..."
	const dotsLen = 3
	ret := input
	if len(input) > size {
		if size > dotsLen {
			size -= dotsLen
		}
		ret = ret[0:size] + dots
	}
	return ret
}

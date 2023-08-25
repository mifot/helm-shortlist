package cmd

import (
	"bytes"
	"fmt"
	"github.com/gosuri/uitable"
	"github.com/spf13/cobra"
	"helm.sh/helm/v3/cmd/helm/require"
	"helm.sh/helm/v3/pkg/action"
	"helm.sh/helm/v3/pkg/cli"
	"helm.sh/helm/v3/pkg/cli/output"
	"helm.sh/helm/v3/pkg/release"
	"io"
	"log"
	"os"
	"time"
)

const shortlistHelp = `
This command shortls lists all of the releases for a specified namespace (uses current namespace context if namespace not specified).
But compare to original helm list command, it prints out only limited information about releases like:
name, namespace, updated time and status. 

By default, it lists only releases that are deployed or failed. Flags like
'--uninstalled' and '--all' will alter this behavior. Such flags can be combined:
'--uninstalled --failed'.

By default, items are sorted alphabetically. Use the '-d' flag to sort by
release date.

If the --filter flag is provided, it will be treated as a filter. Filters are
regular expressions (Perl compatible) that are applied to the list of releases.
Only items that match the filter will be returned.

    $ helm shortls --filter 'te[a-z]+'
    NAME            NAMESPACE       UPDATED                 STATUS  
    another-test    tests           23 Aug 23 16:18 +0200   deployed

If no results are found, 'helm shortls' will exit 0, but with no output.

By default, up to 256 items may be returned. To limit this, use the '--max' flag.
Setting '--max' to 0 will not return all results. Rather, it will return the
server's default, which may be much higher than 256. Pairing the '--max'
flag with the '--offset' flag allows you to page through results.
`

func New() *cobra.Command {

	actionConfig := new(action.Configuration)
	client := action.NewList(actionConfig)
	settings := cli.New()
	requestedNamespace := settings.Namespace()

	cmd := &cobra.Command{
		Use:     "shortls",
		Aliases: []string{"sls"},
		Short:   "Show manifest differences",
		Long:    shortlistHelp,
		Args:    require.NoArgs,

		RunE: func(cmd *cobra.Command, args []string) error {

			out := new(bytes.Buffer)
			outfmt := output.Table

			if client.AllNamespaces {
				if err := actionConfig.Init(settings.RESTClientGetter(), "", os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
					return err
				}
			} else {
				if err := actionConfig.Init(settings.RESTClientGetter(), requestedNamespace, os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
					return err
				}
			}

			client.SetStateMask()

			client.Deployed = true
			results, err := client.Run()
			if err != nil {
				return err
			}

			if err = outfmt.Write(out, newReleaseListWriter(results, client.TimeFormat, client.NoHeaders)); err != nil {
				return err
			}

			fmt.Print(out.String())
			return nil
		},
	}

	cmd.PersistentFlags().BoolVarP(&client.AllNamespaces, "all-namespaces", "A", false, "list releases across all namespaces")
	cmd.PersistentFlags().BoolVarP(&client.NoHeaders, "no-headers", "", false, "don't print headers when using the default output format")
	cmd.PersistentFlags().StringVar(&client.TimeFormat, "time-format", "", `format time using golang time formatter. Example: --time-format "2006-01-02 15:04:05Z0700"`)
	cmd.PersistentFlags().BoolVarP(&client.ByDate, "date", "d", false, "sort by release date")
	cmd.PersistentFlags().BoolVarP(&client.SortReverse, "reverse", "r", false, "reverse the sort order")
	cmd.PersistentFlags().BoolVarP(&client.All, "all", "a", false, "show all releases without any filter applied")
	cmd.PersistentFlags().BoolVar(&client.Uninstalled, "uninstalled", false, "show uninstalled releases (if 'helm uninstall --keep-history' was used)")
	cmd.PersistentFlags().BoolVar(&client.Superseded, "superseded", false, "show superseded releases")
	cmd.PersistentFlags().BoolVar(&client.Uninstalling, "uninstalling", false, "show releases that are currently being uninstalled")
	cmd.PersistentFlags().BoolVar(&client.Deployed, "deployed", false, "show deployed releases. If no other is specified, this will be automatically enabled")
	cmd.PersistentFlags().BoolVar(&client.Failed, "failed", false, "show failed releases")
	cmd.PersistentFlags().BoolVar(&client.Pending, "pending", false, "show pending releases")
	cmd.PersistentFlags().IntVarP(&client.Limit, "max", "m", 256, "maximum number of releases to fetch")
	cmd.PersistentFlags().IntVar(&client.Offset, "offset", 0, "next release index in the list, used to offset from start value")
	cmd.PersistentFlags().StringVarP(&client.Filter, "filter", "f", "", "a regular expression (Perl compatible). Any releases that match the expression will be included in the results")
	cmd.PersistentFlags().StringVarP(&client.Selector, "selector", "l", "", "Selector (label query) to filter on, supports '=', '==', and '!='.(e.g. -l key1=value1,key2=value2). Works only for secret(default) and configmap storage backends.")
	cmd.PersistentFlags().StringVarP(&requestedNamespace, "namespace", "n", settings.Namespace(), "namespace scope for this request")

	cmd.SetFlagErrorFunc(func(command *cobra.Command, err error) error {
		fmt.Printf("%v\n", err)
		return nil
	})

	cmd.SetHelpCommand(&cobra.Command{})
	return cmd
}

type releaseElement struct {
	Name      string `json:"name"`
	Namespace string `json:"namespace"`
	Updated   string `json:"updated"`
	Status    string `json:"status"`
}

type releaseListWriter struct {
	releases  []releaseElement
	noHeaders bool
}

func newReleaseListWriter(releases []*release.Release, timeFormat string, noHeaders bool) *releaseListWriter {
	// Initialize the array so no results returns an empty array instead of null
	elements := make([]releaseElement, 0, len(releases))
	for _, r := range releases {
		element := releaseElement{
			Name:      r.Name,
			Namespace: r.Namespace,
			Status:    r.Info.Status.String(),
		}

		t := "-"
		if tspb := r.Info.LastDeployed; !tspb.IsZero() {
			if timeFormat != "" {
				t = tspb.Format(timeFormat)
			} else {
				t = tspb.Format(time.RFC822Z)
			}
		}

		element.Updated = t

		elements = append(elements, element)
	}
	return &releaseListWriter{elements, noHeaders}
}

func (r *releaseListWriter) WriteTable(out io.Writer) error {
	table := uitable.New()
	if !r.noHeaders {
		table.AddRow("NAME", "NAMESPACE", "UPDATED", "STATUS")
	}
	for _, r := range r.releases {
		table.AddRow(r.Name, r.Namespace, r.Updated, r.Status)
	}
	return output.EncodeTable(out, table)
}

func (r *releaseListWriter) WriteJSON(out io.Writer) error {
	return output.EncodeJSON(out, r.releases)
}

func (r *releaseListWriter) WriteYAML(out io.Writer) error {
	return output.EncodeYAML(out, r.releases)
}

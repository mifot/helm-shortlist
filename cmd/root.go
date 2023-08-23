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

// import (
// 	"fmt"
// // 	"io"
// // 	"os"
// // 	"strconv"

// // 	"github.com/gosuri/uitable"
// // 	"github.com/spf13/cobra"

// // 	"helm.sh/helm/v3/cmd/helm/require"
// // 	"helm.sh/helm/v3/pkg/action"
// // 	"helm.sh/helm/v3/pkg/cli/output"
// // 	"helm.sh/helm/v3/pkg/release"
// )

const listHelp = `
This command lists all of the releases for a specified namespace (uses current namespace context if namespace not specified).

By default, it lists only releases that are deployed or failed. Flags like
'--uninstalled' and '--all' will alter this behavior. Such flags can be combined:
'--uninstalled --failed'.

By default, items are sorted alphabetically. Use the '-d' flag to sort by
release date.

If the --filter flag is provided, it will be treated as a filter. Filters are
regular expressions (Perl compatible) that are applied to the list of releases.
Only items that match the filter will be returned.

    $ helm list --filter 'ara[a-z]+'
    NAME                UPDATED                                  CHART
    maudlin-arachnid    2020-06-18 14:17:46.125134977 +0000 UTC  alpine-0.1.0

If no results are found, 'helm list' will exit 0, but with no output (or in
the case of no '-q' flag, only headers).

By default, up to 256 items may be returned. To limit this, use the '--max' flag.
Setting '--max' to 0 will not return all results. Rather, it will return the
server's default, which may be much higher than 256. Pairing the '--max'
flag with the '--offset' flag allows you to page through results.
`

func New() *cobra.Command {

	cmd := &cobra.Command{
		Use:   "shortls",
		Short: "Show manifest differences",
		Long:  listHelp,
		//Alias root command to chart subcommand
		Args: require.NoArgs,
		// parse the flags and check for actions like suppress-secrets, no-colors
		//PersistentPreRun: func(cmd *cobra.Command, args []string) {
		//	//var fc *bool
		//
		//	//if cmd.Flags().Changed("color") {
		//	//	v, _ := cmd.Flags().GetBool("color")
		//	//	fc = &v
		//	//} else {
		//	//	v, err := strconv.ParseBool(os.Getenv("HELM_DIFF_COLOR"))
		//	//	if err == nil {
		//	//		fc = &v
		//	//	}
		//	//}
		//
		//	//if !cmd.Flags().Changed("output") {
		//	//	v, set := os.LookupEnv("HELM_DIFF_OUTPUT")
		//	//	if set && strings.TrimSpace(v) != "" {
		//	//		_ = cmd.Flags().Set("output", v)
		//	//	}
		//	//}
		//
		//	//nc, _ := cmd.Flags().GetBool("no-color")
		//	//
		//	//if nc || (fc != nil && !*fc) {
		//	//	ansi.DisableColors(true)
		//	//} else if !cmd.Flags().Changed("no-color") && fc == nil {
		//	//	term := term.IsTerminal(int(os.Stdout.Fd()))
		//	//	// https://github.com/databus23/helm-diff/issues/281
		//	//	dumb := os.Getenv("TERM") == "dumb"
		//	//	ansi.DisableColors(!term || dumb)
		//	//}
		//
		//	list()
		//},
		RunE: func(cmd *cobra.Command, args []string) error {
			out := new(bytes.Buffer)
			list(out)
			fmt.Print(out.String())
			//fmt.Print(listOut)
			return nil
		},
	}

	// add no-color as global flag
	cmd.PersistentFlags().Bool("no-color", false, "remove colors from the output. If both --no-color and --color are unspecified, coloring enabled only when the stdout is a term and TERM is not \"dumb\"")
	cmd.PersistentFlags().Bool("color", false, "color output. You can control the value for this flag via HELM_DIFF_COLOR=[true|false]. If both --no-color and --color are unspecified, coloring enabled only when the stdout is a term and TERM is not \"dumb\"")
	// add flagset from chartCommand
	//cmd.Flags().AddFlagSet(chartCommand.Flags())
	//cmd.AddCommand(newVersionCmd(), chartCommand)
	// add subcommands
	//cmd.AddCommand(
	//	revisionCmd(),
	//	rollbackCmd(),
	//	releaseCmd(),
	//)
	cmd.SetHelpCommand(&cobra.Command{}) // Disable the help command
	return cmd
}

func list(out io.Writer) error {
	settings := cli.New()

	//var outfmt output.Format
	outfmt := output.Table // todo improve to support all outputs

	actionConfig := new(action.Configuration)
	// You can pass an empty string instead of settings.Namespace() to list
	// all namespaces
	if err := actionConfig.Init(settings.RESTClientGetter(), settings.Namespace(), os.Getenv("HELM_DRIVER"), log.Printf); err != nil {
		log.Printf("%+v", err)
		os.Exit(1)
	}

	client := action.NewList(actionConfig)
	//client.TimeFormat =
	// Only list deployed
	client.Deployed = true
	results, err := client.Run()
	if err != nil {
		log.Printf("%+v", err)
		os.Exit(1)
		return err
	}

	//for _, rel := range results {
	//	//log.Printf("%+v", rel)
	//}

	return outfmt.Write(out, newReleaseListWriter(results, client.TimeFormat, client.NoHeaders))

	//return newReleaseListWriter(results, client.TimeFormat, client.NoHeaders)
}

// func newListCmd(cfg *action.Configuration, out io.Writer) *cobra.Command {
// 	client := action.NewList(cfg)
// 	var outfmt output.Format

// 	cmd := &cobra.Command{
// 		Use:               "list",
// 		Short:             "list releases",
// 		Long:              listHelp,
// 		Aliases:           []string{"ls"},
// 		Args:              require.NoArgs,
// 		ValidArgsFunction: noCompletions,
// 		RunE: func(cmd *cobra.Command, args []string) error {
// 			if client.AllNamespaces {
// 				if err := cfg.Init(settings.RESTClientGetter(), "", os.Getenv("HELM_DRIVER"), debug); err != nil {
// 					return err
// 				}
// 			}
// 			client.SetStateMask()

// 			results, err := client.Run()
// 			if err != nil {
// 				return err
// 			}

// 			if client.Short {
// 				names := make([]string, 0, len(results))
// 				for _, res := range results {
// 					names = append(names, res.Name)
// 				}

// 				outputFlag := cmd.Flag("output")

// 				switch outputFlag.Value.String() {
// 				case "json":
// 					output.EncodeJSON(out, names)
// 					return nil
// 				case "yaml":
// 					output.EncodeYAML(out, names)
// 					return nil
// 				case "table":
// 					for _, res := range results {
// 						fmt.Fprintln(out, res.Name)
// 					}
// 					return nil
// 				}
// 			}

// 			return outfmt.Write(out, newReleaseListWriter(results, client.TimeFormat, client.NoHeaders))
// 		},
// 	}

// 	f := cmd.Flags()
// 	f.BoolVarP(&client.Short, "short", "q", false, "output short (quiet) listing format")
// 	f.BoolVarP(&client.NoHeaders, "no-headers", "", false, "don't print headers when using the default output format")
// 	f.StringVar(&client.TimeFormat, "time-format", "", `format time using golang time formatter. Example: --time-format "2006-01-02 15:04:05Z0700"`)
// 	f.BoolVarP(&client.ByDate, "date", "d", false, "sort by release date")
// 	f.BoolVarP(&client.SortReverse, "reverse", "r", false, "reverse the sort order")
// 	f.BoolVarP(&client.All, "all", "a", false, "show all releases without any filter applied")
// 	f.BoolVar(&client.Uninstalled, "uninstalled", false, "show uninstalled releases (if 'helm uninstall --keep-history' was used)")
// 	f.BoolVar(&client.Superseded, "superseded", false, "show superseded releases")
// 	f.BoolVar(&client.Uninstalling, "uninstalling", false, "show releases that are currently being uninstalled")
// 	f.BoolVar(&client.Deployed, "deployed", false, "show deployed releases. If no other is specified, this will be automatically enabled")
// 	f.BoolVar(&client.Failed, "failed", false, "show failed releases")
// 	f.BoolVar(&client.Pending, "pending", false, "show pending releases")
// 	f.BoolVarP(&client.AllNamespaces, "all-namespaces", "A", false, "list releases across all namespaces")
// 	f.IntVarP(&client.Limit, "max", "m", 256, "maximum number of releases to fetch")
// 	f.IntVar(&client.Offset, "offset", 0, "next release index in the list, used to offset from start value")
// 	f.StringVarP(&client.Filter, "filter", "f", "", "a regular expression (Perl compatible). Any releases that match the expression will be included in the results")
// 	f.StringVarP(&client.Selector, "selector", "l", "", "Selector (label query) to filter on, supports '=', '==', and '!='.(e.g. -l key1=value1,key2=value2). Works only for secret(default) and configmap storage backends.")
// 	bindOutputFlag(cmd, &outfmt)

// 	return cmd
// }

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
			t = tspb.Format(time.RFC822Z)
			//t = tspb.Format(time.RFC3339)
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

// // Returns all releases from 'releases', except those with names matching 'ignoredReleases'
// func filterReleases(releases []*release.Release, ignoredReleaseNames []string) []*release.Release {
// 	// if ignoredReleaseNames is nil, just return releases
// 	if ignoredReleaseNames == nil {
// 		return releases
// 	}

// 	var filteredReleases []*release.Release
// 	for _, rel := range releases {
// 		found := false
// 		for _, ignoredName := range ignoredReleaseNames {
// 			if rel.Name == ignoredName {
// 				found = true
// 				break
// 			}
// 		}
// 		if !found {
// 			filteredReleases = append(filteredReleases, rel)
// 		}
// 	}

// 	return filteredReleases
// }

// // Provide dynamic auto-completion for release names
// func compListReleases(toComplete string, ignoredReleaseNames []string, cfg *action.Configuration) ([]string, cobra.ShellCompDirective) {
// 	cobra.CompDebugln(fmt.Sprintf("compListReleases with toComplete %s", toComplete), settings.Debug)

// 	client := action.NewList(cfg)
// 	client.All = true
// 	client.Limit = 0
// 	// Do not filter so as to get the entire list of releases.
// 	// This will allow zsh and fish to match completion choices
// 	// on other criteria then prefix.  For example:
// 	//   helm status ingress<TAB>
// 	// can match
// 	//   helm status nginx-ingress
// 	//
// 	// client.Filter = fmt.Sprintf("^%s", toComplete)

// 	client.SetStateMask()
// 	releases, err := client.Run()
// 	if err != nil {
// 		return nil, cobra.ShellCompDirectiveDefault
// 	}

// 	var choices []string
// 	filteredReleases := filterReleases(releases, ignoredReleaseNames)
// 	for _, rel := range filteredReleases {
// 		choices = append(choices,
// 			fmt.Sprintf("%s\t%s-%s -> %s", rel.Name, rel.Chart.Metadata.Name, rel.Chart.Metadata.Version, rel.Info.Status.String()))
// 	}

// 	return choices, cobra.ShellCompDirectiveNoFileComp
// }

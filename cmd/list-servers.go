package cmd

import (
	"fmt"
	"os"
	"strings"
	"text/tabwriter"

	"github.com/gr00by87/fst/config"
	"github.com/gr00by87/fst/core"
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"
)

var (
	name       *[]string
	env        *[]string
	region     *[]string
	ignoreCase *bool

	// listServersCmd represents the list-servers command.
	listServersCmd = &cobra.Command{
		Use:   "list-servers",
		Short: "List available servers",
		Long:  "This subcommand lists available servers from selected AWS region(s), filtered by Name and Env tags",
		Run:   runListServers,
	}
)

// init initializes the cobra command and flags.
func init() {
	rootCmd.AddCommand(listServersCmd)
	name = listServersCmd.Flags().StringSliceP("name", "n", []string{}, "filter servers by Name tag, multiple comma separated values are allowed")
	env = listServersCmd.Flags().StringSliceP("env", "e", []string{}, "filter servers by Env tag, multiple comma separated values are allowed")
	region = listServersCmd.Flags().StringSliceP("region", "r", []string{"us-east-1"}, "look for servers in selected AWS region(s), any of: us-east-1,us-west-2,eu-west-1,ap-northeast-1,ap-southeast-2,all")
	ignoreCase = listServersCmd.Flags().BoolP("ignore-case", "i", false, "ignore case in tag filters")

	// Remove confusing `[]` symbols from region's default value.
	listServersCmd.Flags().VisitAll(func(flag *pflag.Flag) {
		if flag.Name == "region" {
			flag.DefValue = strings.Map(func(r rune) rune {
				if r == '[' || r == ']' {
					return -1
				}
				return r
			}, flag.DefValue)
		}
	})
}

// runListServers executes the list-servers command.
func runListServers(_ *cobra.Command, _ []string) {
	cfg, err := config.LoadFromFile()
	if err != nil {
		exitWithError(err)
	}

	nameFilter := core.NewFilter(*name, core.Contains, *ignoreCase)
	envFilter := core.NewFilter(*env, core.Equals, *ignoreCase)
	servers, err := core.GetAllServers(cfg.AWSCredentials, nameFilter, envFilter, *region)
	if err != nil {
		exitWithError(err)
	}

	w := tabwriter.NewWriter(os.Stdout, 0, 0, 4, ' ', 0)
	fmt.Fprintln(w, "NAME\tENVIRONMENT\tADDRESS")
	for _, server := range servers {
		fmt.Fprintf(w, "%s\t%s\t%s\n", server.Name, server.Env, server.Address)
	}
	w.Flush()
}
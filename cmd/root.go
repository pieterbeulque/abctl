package cmd

import (
	"errors"
	"github.com/airbytehq/abctl/cmd/local"
	"github.com/airbytehq/abctl/cmd/version"
	localcmd "github.com/airbytehq/abctl/internal/local"
	"github.com/pterm/pterm"
	"os"

	"github.com/spf13/cobra"
)

// Help messages to display for specific error situations.
const (
	// helpDocker is displayed if ErrDocker is ever returned
	helpDocker = `An error occurred while communicating with the Docker daemon.
Ensure that Docker is running and is accessible.  You may need to upgrade to a newer version of Docker.
For additional help please visit https://docs.docker.com/get-docker/`

	// helpKubernetes is displayed if ErrKubernetes is ever returned
	helpKubernetes = `An error occurred while communicating with the Kubernetes cluster.
If using Docker Desktop, ensure that Kubernetes is enabled.
For additional help please visit https://docs.docker.com/desktop/kubernetes/`
)

var (
	// flagDNT indicates if the do-not-track flag was specified
	flagDNT bool

	// rootCmd represents the base command when called without any subcommands
	rootCmd = &cobra.Command{
		Use:   "abctl",
		Short: pterm.LightBlue("Airbyte") + "'s command line tool",
		PersistentPreRun: func(cmd *cobra.Command, args []string) {
			if flagDNT {
				pterm.Info.Println("telemetry disabled (--dnt)")
			}
		},
	}
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		pterm.Error.Println(err)

		if errors.Is(err, localcmd.ErrDocker) {
			pterm.Println()
			pterm.Info.Println(helpDocker)
		} else if errors.Is(err, localcmd.ErrKubernetes) {
			pterm.Println()
			pterm.Info.Println(helpKubernetes)
		}
		os.Exit(1)
	}
}

func init() {
	// configure cobra to chain Persistent*Run commands together
	cobra.EnableTraverseRunHooks = true

	rootCmd.SilenceUsage = true
	rootCmd.SilenceErrors = true

	rootCmd.AddCommand(version.Cmd)
	rootCmd.AddCommand(local.Cmd)
	rootCmd.CompletionOptions.DisableDefaultCmd = true
	rootCmd.PersistentFlags().BoolVar(&flagDNT, "dnt", false, "opt out of telemetry data collection")
}

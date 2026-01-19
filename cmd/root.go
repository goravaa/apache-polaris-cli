package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	Version = "dev"
)

var rootCmd = &cobra.Command{
	Use:   "polaris",
	Short: "A CLI for Apache Polaris",
	Long: `Apache Polaris CLI is a command-line interface for managing
Apache Polaris catalogs, namespaces, tables, and more.

To get started, configure your Polaris server and authenticate:
  polaris config set --host http://localhost:8181
  polaris auth login --client-id <your-client-id> --client-secret <your-client-secret>`,
	Version: Version,
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Enable verbose output")
}

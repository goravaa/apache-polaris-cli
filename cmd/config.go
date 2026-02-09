package cmd

import (
	"fmt"

	"github.com/goravaa/apache-polaris-cli/pkg/config"
	"github.com/spf13/cobra"
)

var (
	configHost          string
	configRealm         string
	configCatalogPrefix string
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration commands",
	Long:  `Commands for managing CLI configuration.`,
}

var configSetCmd = &cobra.Command{
	Use:   "set",
	Short: "Set configuration values",
	Long: `Set configuration values for the Polaris CLI.

Examples:
  polaris config set --host http://localhost:8181
  polaris config set --host https://polaris.example.com --realm my-realm`,
	RunE: runConfigSet,
}

var configShowCmd = &cobra.Command{
	Use:   "show",
	Short: "Show current configuration",
	Long:  `Display the current CLI configuration.`,
	RunE:  runConfigShow,
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.AddCommand(configSetCmd)
	configCmd.AddCommand(configShowCmd)

	configSetCmd.Flags().StringVar(&configHost, "host", "", "Polaris server URL (e.g., http://localhost:8181)")
	configSetCmd.Flags().StringVar(&configRealm, "realm", "", "Polaris realm (for multi-tenant setups)")
	configSetCmd.Flags().StringVar(&configCatalogPrefix, "catalog-prefix", "", "Default catalog prefix for catalog API calls")
}

func runConfigSet(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		cfg = &config.Config{}
	}

	if configHost != "" {
		cfg.Host = configHost
	}
	if cmd.Flags().Changed("realm") {
		cfg.Realm = configRealm
	}
	if cmd.Flags().Changed("catalog-prefix") {
		cfg.CatalogPrefix = configCatalogPrefix
	}

	if cfg.Host == "" {
		return fmt.Errorf("host is required. Use --host to set the Polaris server URL")
	}

	if err := config.SaveConfig(cfg); err != nil {
		return fmt.Errorf("failed to save config: %w", err)
	}

	fmt.Println("✓ Configuration saved!")
	fmt.Printf("  Host: %s\n", cfg.Host)
	if cfg.Realm != "" {
		fmt.Printf("  Realm: %s\n", cfg.Realm)
	}
	if cfg.CatalogPrefix != "" {
		fmt.Printf("  Catalog Prefix: %s\n", cfg.CatalogPrefix)
	}

	return nil
}

func runConfigShow(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		return fmt.Errorf("failed to load config: %w", err)
	}

	fmt.Println("Current Configuration:")
	fmt.Printf("  Host: %s\n", cfg.Host)
	if cfg.Realm != "" {
		fmt.Printf("  Realm: %s\n", cfg.Realm)
	} else {
		fmt.Println("  Realm: (not set)")
	}
	if cfg.CatalogPrefix != "" {
		fmt.Printf("  Catalog Prefix: %s\n", cfg.CatalogPrefix)
	} else {
		fmt.Println("  Catalog Prefix: (not set)")
	}

	return nil
}

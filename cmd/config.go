package cmd

import (
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/goravaa/apache-polaris-cli/pkg/config"
	"github.com/spf13/cobra"
)

var (
	configHost          string
	configRealm         string
	configCatalogPrefix string
	configRootClientID  string
	configRootSecret    string
	configFile          string
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Configuration commands",
	Long:  `Commands for managing CLI configuration.`,
}

var configSetCmd = &cobra.Command{
	Use:   "set [config-json-or-yaml]",
	Short: "Set configuration values",
	Long: `Set configuration values for the Polaris CLI.

Examples:
  polaris config set --host http://localhost:8181
  polaris config set --host https://polaris.example.com --realm my-realm
  polaris config set --file ./polaris.yaml
  polaris config set '{"host":"http://localhost:8181","realm":"default-realm","root_client_id":"root","root_client_secret":"secret"}'
  polaris config set - < ./polaris.yaml`,
	Args: cobra.MaximumNArgs(1),
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
	configSetCmd.Flags().StringVar(&configRootClientID, "root-client-id", "", "Default root OAuth client ID for login")
	configSetCmd.Flags().StringVar(&configRootSecret, "root-client-secret", "", "Default root OAuth client secret for login")
	configSetCmd.Flags().StringVarP(&configFile, "file", "f", "", "Read JSON or YAML config from a file")
}

func runConfigSet(cmd *cobra.Command, args []string) error {
	cfg, err := config.LoadConfig()
	if err != nil {
		cfg = &config.Config{}
	}

	if configFile != "" {
		values, err := readConfigValuesFromFile(configFile)
		if err != nil {
			return err
		}
		if err := config.ApplyConfigValues(cfg, values); err != nil {
			return err
		}
	}

	if len(args) == 1 {
		values, err := readConfigValuesFromArg(args[0])
		if err != nil {
			return err
		}
		if err := config.ApplyConfigValues(cfg, values); err != nil {
			return err
		}
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
	if cmd.Flags().Changed("root-client-id") {
		cfg.RootClientID = configRootClientID
	}
	if cmd.Flags().Changed("root-client-secret") {
		cfg.RootClientSecret = configRootSecret
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
	if cfg.RootClientID != "" {
		fmt.Printf("  Root Client ID: %s\n", cfg.RootClientID)
	}
	if cfg.RootClientSecret != "" {
		fmt.Println("  Root Client Secret: (set)")
	}

	return nil
}

func readConfigValuesFromFile(path string) (map[string]string, error) {
	var data []byte
	var err error
	if path == "-" {
		data, err = io.ReadAll(os.Stdin)
	} else {
		data, err = os.ReadFile(path)
	}
	if err != nil {
		return nil, fmt.Errorf("failed to read config input: %w", err)
	}
	return parseConfigValues(data)
}

func readConfigValuesFromArg(arg string) (map[string]string, error) {
	if arg == "-" {
		return readConfigValuesFromFile(arg)
	}

	if looksLikeInlineConfig(arg) {
		return parseConfigValues([]byte(arg))
	}

	data, err := os.ReadFile(arg)
	if err != nil {
		return nil, fmt.Errorf("failed to read config argument as a file: %w", err)
	}
	return parseConfigValues(data)
}

func parseConfigValues(data []byte) (map[string]string, error) {
	values, err := config.ParseConfigValues(data)
	if err != nil {
		return nil, err
	}
	return values, nil
}

func looksLikeInlineConfig(input string) bool {
	if strings.Contains(input, "\n") || strings.Contains(input, ": ") {
		return true
	}

	for _, r := range input {
		switch r {
		case ' ', '\t', '\n', '\r':
			continue
		case '{', '[':
			return true
		default:
			return false
		}
	}
	return false
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
	if cfg.RootClientID != "" {
		fmt.Printf("  Root Client ID: %s\n", cfg.RootClientID)
	} else {
		fmt.Println("  Root Client ID: (not set)")
	}
	if cfg.RootClientSecret != "" {
		fmt.Println("  Root Client Secret: (set)")
	} else {
		fmt.Println("  Root Client Secret: (not set)")
	}

	return nil
}

package cmd

import (
	"context"
	"fmt"

	catalogapi "github.com/goravaa/apache-polaris-cli/pkg/api/openapi/catalog"
	"github.com/spf13/cobra"
)

var tableNamespace string

var catalogTablesCmd = &cobra.Command{
	Use:   "tables",
	Short: "Table operations",
}

var catalogTablesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List tables in a namespace",
	RunE:  runCatalogTablesList,
}

func init() {
	catalogCmd.AddCommand(catalogTablesCmd)
	catalogTablesCmd.AddCommand(catalogTablesListCmd)

	catalogTablesListCmd.Flags().StringVar(&tableNamespace, "namespace", "", "Namespace (dot- or slash-separated)")
}

func runCatalogTablesList(cmd *cobra.Command, args []string) error {
	if tableNamespace == "" {
		return fmt.Errorf("--namespace is required")
	}

	client, cfg, err := newCatalogClient()
	if err != nil {
		return err
	}

	prefix, err := resolveCatalogPrefix(cfg)
	if err != nil {
		return err
	}

	parts, err := parseNamespaceArg(tableNamespace)
	if err != nil {
		return err
	}

	nsPath := namespacePath(parts)
	resp, err := client.ListTablesWithResponse(
		context.Background(),
		catalogapi.Prefix(prefix),
		catalogapi.NamespaceString(nsPath),
		nil,
	)
	if err != nil {
		return err
	}
	if resp.JSON200 == nil {
		return fmt.Errorf("request failed: %s", resp.Status())
	}

	if resp.JSON200.Identifiers == nil || len(*resp.JSON200.Identifiers) == 0 {
		fmt.Println("(no tables)")
		return nil
	}

	for _, ident := range *resp.JSON200.Identifiers {
		fmt.Printf("%s.%s\n", formatNamespace(ident.Namespace), ident.Name)
	}

	return nil
}

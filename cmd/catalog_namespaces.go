package cmd

import (
	"context"
	"fmt"

	catalogapi "github.com/goravaa/apache-polaris-cli/pkg/api/openapi/catalog"
	"github.com/spf13/cobra"
)

var namespaceProperties []string

var catalogNamespacesCmd = &cobra.Command{
	Use:   "namespaces",
	Short: "Namespace operations",
}

var catalogNamespacesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List namespaces",
	RunE:  runCatalogNamespacesList,
}

var catalogNamespacesCreateCmd = &cobra.Command{
	Use:   "create <namespace>",
	Short: "Create a namespace",
	Args:  cobra.ExactArgs(1),
	RunE:  runCatalogNamespacesCreate,
}

var catalogNamespacesDeleteCmd = &cobra.Command{
	Use:   "delete <namespace>",
	Short: "Delete an empty namespace",
	Args:  cobra.ExactArgs(1),
	RunE:  runCatalogNamespacesDelete,
}

func init() {
	catalogCmd.AddCommand(catalogNamespacesCmd)
	catalogNamespacesCmd.AddCommand(catalogNamespacesListCmd)
	catalogNamespacesCmd.AddCommand(catalogNamespacesCreateCmd)
	catalogNamespacesCmd.AddCommand(catalogNamespacesDeleteCmd)

	catalogNamespacesCreateCmd.Flags().StringArrayVar(&namespaceProperties, "property", nil, "Namespace property key=value (repeatable)")
}

func runCatalogNamespacesList(cmd *cobra.Command, args []string) error {
	client, cfg, err := newCatalogClient()
	if err != nil {
		return err
	}

	prefix, err := resolveCatalogPrefix(cfg)
	if err != nil {
		return err
	}

	resp, err := client.ListNamespacesWithResponse(context.Background(), catalogapi.Prefix(prefix), nil)
	if err != nil {
		return err
	}
	if resp.JSON200 == nil {
		return fmt.Errorf("request failed: %s", resp.Status())
	}

	if resp.JSON200.Namespaces == nil || len(*resp.JSON200.Namespaces) == 0 {
		fmt.Println("(no namespaces)")
		return nil
	}

	for _, ns := range *resp.JSON200.Namespaces {
		fmt.Println(formatNamespace(ns))
	}

	return nil
}

func runCatalogNamespacesCreate(cmd *cobra.Command, args []string) error {
	client, cfg, err := newCatalogClient()
	if err != nil {
		return err
	}

	prefix, err := resolveCatalogPrefix(cfg)
	if err != nil {
		return err
	}

	parts, err := parseNamespaceArg(args[0])
	if err != nil {
		return err
	}

	props, err := parseProperties(namespaceProperties)
	if err != nil {
		return err
	}

	req := catalogapi.CreateNamespaceRequest{
		Namespace: catalogapi.Namespace(parts),
	}
	if len(props) > 0 {
		req.Properties = &props
	}

	resp, err := client.CreateNamespaceWithResponse(context.Background(), catalogapi.Prefix(prefix), req)
	if err != nil {
		return err
	}
	if resp.JSON200 == nil {
		return fmt.Errorf("request failed: %s", resp.Status())
	}

	fmt.Printf("Created namespace %s\n", formatNamespace(parts))
	return nil
}

func runCatalogNamespacesDelete(cmd *cobra.Command, args []string) error {
	client, cfg, err := newCatalogClient()
	if err != nil {
		return err
	}

	prefix, err := resolveCatalogPrefix(cfg)
	if err != nil {
		return err
	}

	parts, err := parseNamespaceArg(args[0])
	if err != nil {
		return err
	}

	resp, err := client.DropNamespaceWithResponse(
		context.Background(),
		catalogapi.Prefix(prefix),
		catalogapi.NamespaceString(namespacePath(parts)),
	)
	if err != nil {
		return err
	}
	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		return fmt.Errorf("request failed: %s", resp.Status())
	}

	fmt.Printf("Deleted namespace %s\n", formatNamespace(parts))
	return nil
}

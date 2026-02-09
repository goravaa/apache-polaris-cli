package cmd

import (
	"context"
	"fmt"
	"strings"

	managementapi "github.com/goravaa/apache-polaris-cli/pkg/api/openapi/management"
	"github.com/spf13/cobra"
)

var (
	catalogName                string
	catalogType                string
	catalogStorageType         string
	catalogDefaultBaseLocation string
	catalogAllowedLocations    []string
	catalogProperties          []string
)

var catalogsCmd = &cobra.Command{
	Use:   "catalogs",
	Short: "Catalog management commands",
}

var catalogsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List catalogs",
	RunE:  runCatalogsList,
}

var catalogsDescribeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Describe a catalog",
	RunE:  runCatalogsDescribe,
}

var catalogsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a catalog",
	RunE:  runCatalogsCreate,
}

var catalogsDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a catalog",
	RunE:  runCatalogsDelete,
}

func init() {
	rootCmd.AddCommand(catalogsCmd)
	catalogsCmd.AddCommand(catalogsListCmd)
	catalogsCmd.AddCommand(catalogsDescribeCmd)
	catalogsCmd.AddCommand(catalogsCreateCmd)
	catalogsCmd.AddCommand(catalogsDeleteCmd)

	catalogsDescribeCmd.Flags().StringVar(&catalogName, "name", "", "Catalog name (required)")

	catalogsCreateCmd.Flags().StringVar(&catalogName, "name", "", "Catalog name (required)")
	catalogsCreateCmd.Flags().StringVar(&catalogType, "type", "INTERNAL", "Catalog type: INTERNAL or EXTERNAL")
	catalogsCreateCmd.Flags().StringVar(&catalogStorageType, "storage-type", "S3", "Storage type: S3, GCS, AZURE, FILE")
	catalogsCreateCmd.Flags().StringVar(&catalogDefaultBaseLocation, "default-base-location", "", "Default base location (required)")
	catalogsCreateCmd.Flags().StringArrayVar(&catalogAllowedLocations, "allowed-location", nil, "Allowed location (repeatable)")
	catalogsCreateCmd.Flags().StringArrayVar(&catalogProperties, "property", nil, "Catalog property key=value (repeatable)")

	catalogsDeleteCmd.Flags().StringVar(&catalogName, "name", "", "Catalog name (required)")
}

func runCatalogsList(cmd *cobra.Command, args []string) error {
	client, _, err := newManagementClient()
	if err != nil {
		return err
	}

	resp, err := client.ListCatalogsWithResponse(context.Background())
	if err != nil {
		return err
	}
	if resp.JSON200 == nil {
		return fmt.Errorf("request failed: %s", resp.Status())
	}

	if len(resp.JSON200.Catalogs) == 0 {
		fmt.Println("(no catalogs)")
		return nil
	}

	for _, c := range resp.JSON200.Catalogs {
		fmt.Println(c.Name)
	}

	return nil
}

func runCatalogsCreate(cmd *cobra.Command, args []string) error {
	if catalogName == "" {
		return fmt.Errorf("--name is required")
	}
	if catalogDefaultBaseLocation == "" {
		return fmt.Errorf("--default-base-location is required")
	}

	typ, err := parseCatalogType(catalogType)
	if err != nil {
		return err
	}

	storageType, err := parseStorageType(catalogStorageType)
	if err != nil {
		return err
	}

	props, err := parseProperties(catalogProperties)
	if err != nil {
		return err
	}

	catalogProps := managementapi.Catalog_Properties{
		DefaultBaseLocation: catalogDefaultBaseLocation,
	}
	if len(props) > 0 {
		catalogProps.AdditionalProperties = props
	}

	storage := managementapi.StorageConfigInfo{
		StorageType: storageType,
	}
	if len(catalogAllowedLocations) > 0 {
		storage.AllowedLocations = &catalogAllowedLocations
	}

	req := managementapi.CreateCatalogRequest{
		Catalog: managementapi.Catalog{
			Name:              catalogName,
			Type:              typ,
			Properties:        catalogProps,
			StorageConfigInfo: storage,
		},
	}

	client, _, err := newManagementClient()
	if err != nil {
		return err
	}

	resp, err := client.CreateCatalogWithResponse(context.Background(), req)
	if err != nil {
		return err
	}
	if resp.JSON201 == nil {
		return fmt.Errorf("request failed: %s", resp.Status())
	}

	fmt.Printf("Created catalog %s\n", catalogName)
	return nil
}

func runCatalogsDescribe(cmd *cobra.Command, args []string) error {
	if catalogName == "" {
		return fmt.Errorf("--name is required")
	}

	client, _, err := newManagementClient()
	if err != nil {
		return err
	}

	resp, err := client.GetCatalogWithResponse(context.Background(), catalogName)
	if err != nil {
		return err
	}
	if resp.JSON200 == nil {
		return fmt.Errorf("request failed: %s", resp.Status())
	}

	c := resp.JSON200
	fmt.Printf("Name: %s\n", c.Name)
	fmt.Printf("Type: %s\n", c.Type)
	fmt.Printf("Storage Type: %s\n", c.StorageConfigInfo.StorageType)
	if c.Properties.DefaultBaseLocation != "" {
		fmt.Printf("Default Base Location: %s\n", c.Properties.DefaultBaseLocation)
	}
	if c.StorageConfigInfo.AllowedLocations != nil && len(*c.StorageConfigInfo.AllowedLocations) > 0 {
		fmt.Println("Allowed Locations:")
		for _, loc := range *c.StorageConfigInfo.AllowedLocations {
			fmt.Printf("  %s\n", loc)
		}
	}
	if len(c.Properties.AdditionalProperties) > 0 {
		fmt.Println("Properties:")
		for k, v := range c.Properties.AdditionalProperties {
			fmt.Printf("  %s=%s\n", k, v)
		}
	}

	return nil
}

func runCatalogsDelete(cmd *cobra.Command, args []string) error {
	if catalogName == "" {
		return fmt.Errorf("--name is required")
	}

	client, _, err := newManagementClient()
	if err != nil {
		return err
	}

	resp, err := client.DeleteCatalogWithResponse(context.Background(), catalogName)
	if err != nil {
		return err
	}
	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		return fmt.Errorf("request failed: %s", resp.Status())
	}

	fmt.Printf("Deleted catalog %s\n", catalogName)
	return nil
}

func parseCatalogType(input string) (managementapi.CatalogType, error) {
	switch strings.ToUpper(strings.TrimSpace(input)) {
	case "INTERNAL":
		return managementapi.INTERNAL, nil
	case "EXTERNAL":
		return managementapi.EXTERNAL, nil
	default:
		return "", fmt.Errorf("invalid catalog type %q (expected INTERNAL or EXTERNAL)", input)
	}
}

func parseStorageType(input string) (managementapi.StorageConfigInfoStorageType, error) {
	switch strings.ToUpper(strings.TrimSpace(input)) {
	case "S3":
		return managementapi.S3, nil
	case "GCS":
		return managementapi.GCS, nil
	case "AZURE":
		return managementapi.AZURE, nil
	case "FILE":
		return managementapi.FILE, nil
	default:
		return "", fmt.Errorf("invalid storage type %q (expected S3, GCS, AZURE, FILE)", input)
	}
}

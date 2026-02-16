package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"strings"

	managementapi "github.com/goravaa/apache-polaris-cli/pkg/api/openapi/management"
	"github.com/spf13/cobra"
)

var (
	catalogRoleName       string
	catalogRoleProperties []string
)

// Grant related flags
var (
	crGrantType      string
	crGrantPrivilege string
	crGrantNamespace string
	crGrantTable     string
	crGrantView      string
	crRevokeCascade  bool
)

var catalogRolesCmd = &cobra.Command{
	Use:   "roles",
	Short: "Catalog role management commands",
}

var catalogRolesListCmd = &cobra.Command{
	Use:   "list",
	Short: "List catalog roles",
	RunE:  runCatalogRolesList,
}

var catalogRolesCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a catalog role",
	RunE:  runCatalogRolesCreate,
}

var catalogRolesDescribeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Describe a catalog role",
	RunE:  runCatalogRolesDescribe,
}

var catalogRolesDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a catalog role",
	RunE:  runCatalogRolesDelete,
}

var catalogRolesUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a catalog role",
	RunE:  runCatalogRolesUpdate,
}

var catalogRolesGrantCmd = &cobra.Command{
	Use:   "grant",
	Short: "Grant privileges to a catalog role",
	RunE:  runCatalogRolesGrant,
}

var catalogRolesRevokeCmd = &cobra.Command{
	Use:   "revoke",
	Short: "Revoke privileges from a catalog role",
	RunE:  runCatalogRolesRevoke,
}

var catalogRolesGrantsCmd = &cobra.Command{
	Use:   "grants",
	Short: "Manage grants for a catalog role",
}

var catalogRolesGrantsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List grants for a catalog role",
	RunE:  runCatalogRolesGrantsList,
}

var catalogRolesPrincipalsCmd = &cobra.Command{
	Use:   "principals",
	Short: "Manage principals assigned to a catalog role",
}

var catalogRolesPrincipalsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List principal roles assigned to a catalog role",
	RunE:  runCatalogRolesPrincipalsList,
}

func init() {
	catalogsCmd.AddCommand(catalogRolesCmd)
	catalogRolesCmd.AddCommand(catalogRolesListCmd)
	catalogRolesCmd.AddCommand(catalogRolesCreateCmd)
	catalogRolesCmd.AddCommand(catalogRolesDescribeCmd)
	catalogRolesCmd.AddCommand(catalogRolesDeleteCmd)
	catalogRolesCmd.AddCommand(catalogRolesUpdateCmd)
	catalogRolesCmd.AddCommand(catalogRolesGrantCmd)
	catalogRolesCmd.AddCommand(catalogRolesRevokeCmd)
	catalogRolesCmd.AddCommand(catalogRolesGrantsCmd)
	catalogRolesGrantsCmd.AddCommand(catalogRolesGrantsListCmd)
	catalogRolesCmd.AddCommand(catalogRolesPrincipalsCmd)
	catalogRolesPrincipalsCmd.AddCommand(catalogRolesPrincipalsListCmd)

	// Flags for list
	catalogRolesListCmd.Flags().StringVar(&catalogName, "catalog", "", "Catalog name (required)")
	catalogRolesListCmd.MarkFlagRequired("catalog")

	// Flags for create
	catalogRolesCreateCmd.Flags().StringVar(&catalogName, "catalog", "", "Catalog name (required)")
	catalogRolesCreateCmd.Flags().StringVar(&catalogRoleName, "name", "", "Role name (required)")
	catalogRolesCreateCmd.Flags().StringArrayVar(&catalogRoleProperties, "property", nil, "Role property key=value (repeatable)")
	catalogRolesCreateCmd.MarkFlagRequired("catalog")
	catalogRolesCreateCmd.MarkFlagRequired("name")

	// Flags for describe
	catalogRolesDescribeCmd.Flags().StringVar(&catalogName, "catalog", "", "Catalog name (required)")
	catalogRolesDescribeCmd.Flags().StringVar(&catalogRoleName, "name", "", "Role name (required)")
	catalogRolesDescribeCmd.MarkFlagRequired("catalog")
	catalogRolesDescribeCmd.MarkFlagRequired("name")

	// Flags for delete
	catalogRolesDeleteCmd.Flags().StringVar(&catalogName, "catalog", "", "Catalog name (required)")
	catalogRolesDeleteCmd.Flags().StringVar(&catalogRoleName, "name", "", "Role name (required)")
	catalogRolesDeleteCmd.MarkFlagRequired("catalog")
	catalogRolesDeleteCmd.MarkFlagRequired("name")

	// Flags for update
	catalogRolesUpdateCmd.Flags().StringVar(&catalogName, "catalog", "", "Catalog name (required)")
	catalogRolesUpdateCmd.Flags().StringVar(&catalogRoleName, "name", "", "Role name (required)")
	catalogRolesUpdateCmd.Flags().StringArrayVar(&catalogRoleProperties, "property", nil, "Role property key=value (repeatable)")
	catalogRolesUpdateCmd.MarkFlagRequired("catalog")
	catalogRolesUpdateCmd.MarkFlagRequired("name")

	// Flags for grant
	catalogRolesGrantCmd.Flags().StringVar(&catalogName, "catalog", "", "Catalog name (required)")
	catalogRolesGrantCmd.Flags().StringVar(&catalogRoleName, "role", "", "Role name (required)")
	catalogRolesGrantCmd.Flags().StringVar(&crGrantType, "type", "", "Resource type: catalog, namespace, table, view (required)")
	catalogRolesGrantCmd.Flags().StringVar(&crGrantPrivilege, "privilege", "", "Privilege to grant (required)")
	catalogRolesGrantCmd.Flags().StringVar(&crGrantNamespace, "namespace", "", "Namespace (for namespace, table, view grants)")
	catalogRolesGrantCmd.Flags().StringVar(&crGrantTable, "table", "", "Table name (for table grants)")
	catalogRolesGrantCmd.Flags().StringVar(&crGrantView, "view", "", "View name (for view grants)")
	catalogRolesGrantCmd.MarkFlagRequired("catalog")
	catalogRolesGrantCmd.MarkFlagRequired("role")
	catalogRolesGrantCmd.MarkFlagRequired("type")
	catalogRolesGrantCmd.MarkFlagRequired("privilege")

	// Flags for revoke
	catalogRolesRevokeCmd.Flags().StringVar(&catalogName, "catalog", "", "Catalog name (required)")
	catalogRolesRevokeCmd.Flags().StringVar(&catalogRoleName, "role", "", "Role name (required)")
	catalogRolesRevokeCmd.Flags().StringVar(&crGrantType, "type", "", "Resource type: catalog, namespace, table, view (required)")
	catalogRolesRevokeCmd.Flags().StringVar(&crGrantPrivilege, "privilege", "", "Privilege to revoke (required)")
	catalogRolesRevokeCmd.Flags().StringVar(&crGrantNamespace, "namespace", "", "Namespace (for namespace, table, view grants)")
	catalogRolesRevokeCmd.Flags().StringVar(&crGrantTable, "table", "", "Table name (for table grants)")
	catalogRolesRevokeCmd.Flags().StringVar(&crGrantView, "view", "", "View name (for view grants)")
	catalogRolesRevokeCmd.Flags().BoolVar(&crRevokeCascade, "cascade", false, "Cascade revocation")
	catalogRolesRevokeCmd.MarkFlagRequired("catalog")
	catalogRolesRevokeCmd.MarkFlagRequired("role")
	catalogRolesRevokeCmd.MarkFlagRequired("type")
	catalogRolesRevokeCmd.MarkFlagRequired("privilege")

	// Flags for list grants
	catalogRolesGrantsListCmd.Flags().StringVar(&catalogName, "catalog", "", "Catalog name (required)")
	catalogRolesGrantsListCmd.Flags().StringVar(&catalogRoleName, "role", "", "Role name (required)")
	catalogRolesGrantsListCmd.MarkFlagRequired("catalog")
	catalogRolesGrantsListCmd.MarkFlagRequired("role")

	// Flags for list principals
	catalogRolesPrincipalsListCmd.Flags().StringVar(&catalogName, "catalog", "", "Catalog name (required)")
	catalogRolesPrincipalsListCmd.Flags().StringVar(&catalogRoleName, "role", "", "Role name (required)")
	catalogRolesPrincipalsListCmd.MarkFlagRequired("catalog")
	catalogRolesPrincipalsListCmd.MarkFlagRequired("role")
}

func runCatalogRolesList(cmd *cobra.Command, args []string) error {
	client, _, err := newManagementClient()
	if err != nil {
		return err
	}

	resp, err := client.ListCatalogRolesWithResponse(context.Background(), catalogName)
	if err != nil {
		return err
	}
	if resp.JSON200 == nil {
		return fmt.Errorf("request failed: %s", resp.Status())
	}

	if len(resp.JSON200.Roles) == 0 {
		fmt.Println("(no roles)")
		return nil
	}

	for _, r := range resp.JSON200.Roles {
		fmt.Println(r.Name)
	}

	return nil
}

func runCatalogRolesCreate(cmd *cobra.Command, args []string) error {
	props, err := parseProperties(catalogRoleProperties)
	if err != nil {
		return err
	}

	req := managementapi.CreateCatalogRoleRequest{
		CatalogRole: &managementapi.CatalogRole{
			Name: catalogRoleName,
		},
	}
	if len(props) > 0 {
		req.CatalogRole.Properties = &props
	}

	client, _, err := newManagementClient()
	if err != nil {
		return err
	}

	resp, err := client.CreateCatalogRoleWithResponse(context.Background(), catalogName, req)
	if err != nil {
		return err
	}
	if resp.JSON201 == nil {
		return fmt.Errorf("request failed: %s", resp.Status())
	}

	fmt.Printf("Created catalog role %s in catalog %s\n", catalogRoleName, catalogName)
	return nil
}

func runCatalogRolesDescribe(cmd *cobra.Command, args []string) error {
	client, _, err := newManagementClient()
	if err != nil {
		return err
	}

	resp, err := client.GetCatalogRoleWithResponse(context.Background(), catalogName, catalogRoleName)
	if err != nil {
		return err
	}
	if resp.JSON200 == nil {
		return fmt.Errorf("request failed: %s", resp.Status())
	}

	r := resp.JSON200
	fmt.Printf("Name: %s\n", r.Name)
	if r.Properties != nil && len(*r.Properties) > 0 {
		fmt.Println("Properties:")
		for k, v := range *r.Properties {
			fmt.Printf("  %s=%s\n", k, v)
		}
	}
	if r.CreateTimestamp != nil {
		fmt.Printf("Created: %d\n", *r.CreateTimestamp)
	}
	if r.LastUpdateTimestamp != nil {
		fmt.Printf("Last Updated: %d\n", *r.LastUpdateTimestamp)
	}
	if r.EntityVersion != nil {
		fmt.Printf("Version: %d\n", *r.EntityVersion)
	}

	return nil
}

func runCatalogRolesDelete(cmd *cobra.Command, args []string) error {
	client, _, err := newManagementClient()
	if err != nil {
		return err
	}

	resp, err := client.DeleteCatalogRoleWithResponse(context.Background(), catalogName, catalogRoleName)
	if err != nil {
		return err
	}
	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		return fmt.Errorf("request failed: %s", resp.Status())
	}

	fmt.Printf("Deleted catalog role %s from catalog %s\n", catalogRoleName, catalogName)
	return nil
}

func runCatalogRolesUpdate(cmd *cobra.Command, args []string) error {
	client, _, err := newManagementClient()
	if err != nil {
		return err
	}

	// Fetch current role to get version
	roleResp, err := client.GetCatalogRoleWithResponse(context.Background(), catalogName, catalogRoleName)
	if err != nil {
		return err
	}
	if roleResp.JSON200 == nil {
		return fmt.Errorf("failed to fetch role: %s", roleResp.Status())
	}

	currentVersion := 0
	if roleResp.JSON200.EntityVersion != nil {
		currentVersion = *roleResp.JSON200.EntityVersion
	}

	props, err := parseProperties(catalogRoleProperties)
	if err != nil {
		return err
	}

	req := managementapi.UpdateCatalogRoleRequest{
		CurrentEntityVersion: currentVersion,
		Properties:           props,
	}

	resp, err := client.UpdateCatalogRoleWithResponse(context.Background(), catalogName, catalogRoleName, req)
	if err != nil {
		return err
	}
	if resp.JSON200 == nil {
		return fmt.Errorf("request failed: %s", resp.Status())
	}

	fmt.Printf("Updated catalog role %s in catalog %s\n", catalogRoleName, catalogName)
	return nil
}

// Custom Grant Structs (since generated ones are missing)
type CatalogGrant struct {
	Type      string `json:"type"`
	Privilege string `json:"privilege"`
}

type NamespaceGrant struct {
	Type      string   `json:"type"`
	Namespace []string `json:"namespace"`
	Privilege string   `json:"privilege"`
}

type TableGrant struct {
	Type      string   `json:"type"`
	Namespace []string `json:"namespace"`
	TableName string   `json:"tableName"`
	Privilege string   `json:"privilege"`
}

type ViewGrant struct {
	Type      string   `json:"type"`
	Namespace []string `json:"namespace"`
	ViewName  string   `json:"viewName"`
	Privilege string   `json:"privilege"`
}

type AddGrantRequest struct {
	Grant interface{} `json:"grant"`
}

type RevokeGrantRequest struct {
	Grant interface{} `json:"grant"`
}

func constructGrant(typ, priv, nsStr, table, view string) (interface{}, error) {
	typ = strings.ToLower(typ)
	priv = strings.ToUpper(priv)

	var ns []string
	if nsStr != "" {
		var err error
		ns, err = parseNamespaceArg(nsStr)
		if err != nil {
			return nil, err
		}
	}

	switch typ {
	case "catalog":
		return CatalogGrant{
			Type:      "catalog",
			Privilege: priv,
		}, nil
	case "namespace":
		if len(ns) == 0 {
			return nil, fmt.Errorf("namespace is required for namespace grant")
		}
		return NamespaceGrant{
			Type:      "namespace",
			Namespace: ns,
			Privilege: priv,
		}, nil
	case "table":
		if len(ns) == 0 {
			return nil, fmt.Errorf("namespace is required for table grant")
		}
		if table == "" {
			return nil, fmt.Errorf("table is required for table grant")
		}
		return TableGrant{
			Type:      "table",
			Namespace: ns,
			TableName: table,
			Privilege: priv,
		}, nil
	case "view":
		if len(ns) == 0 {
			return nil, fmt.Errorf("namespace is required for view grant")
		}
		if view == "" {
			return nil, fmt.Errorf("view is required for view grant")
		}
		return ViewGrant{
			Type:      "view",
			Namespace: ns,
			ViewName:  view,
			Privilege: priv,
		}, nil
	default:
		return nil, fmt.Errorf("invalid grant type: %s (expected catalog, namespace, table, view)", typ)
	}
}

func runCatalogRolesGrant(cmd *cobra.Command, args []string) error {
	grantObj, err := constructGrant(crGrantType, crGrantPrivilege, crGrantNamespace, crGrantTable, crGrantView)
	if err != nil {
		return err
	}

	req := AddGrantRequest{Grant: grantObj}
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	client, _, err := newManagementClient()
	if err != nil {
		return err
	}

	resp, err := client.AddGrantToCatalogRoleWithBodyWithResponse(context.Background(), catalogName, catalogRoleName, "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		return err
	}
	if resp.StatusCode() != 201 {
		return fmt.Errorf("request failed: %s", resp.Status())
	}

	fmt.Printf("Granted %s on %s to role %s\n", crGrantPrivilege, crGrantType, catalogRoleName)
	return nil
}

func runCatalogRolesRevoke(cmd *cobra.Command, args []string) error {
	grantObj, err := constructGrant(crGrantType, crGrantPrivilege, crGrantNamespace, crGrantTable, crGrantView)
	if err != nil {
		return err
	}

	req := RevokeGrantRequest{Grant: grantObj}
	bodyBytes, err := json.Marshal(req)
	if err != nil {
		return fmt.Errorf("failed to marshal request: %w", err)
	}

	client, _, err := newManagementClient()
	if err != nil {
		return err
	}

	params := &managementapi.RevokeGrantFromCatalogRoleParams{
		Cascade: &crRevokeCascade,
	}

	resp, err := client.RevokeGrantFromCatalogRoleWithBodyWithResponse(context.Background(), catalogName, catalogRoleName, params, "application/json", bytes.NewReader(bodyBytes))
	if err != nil {
		return err
	}
	if resp.StatusCode() != 204 && resp.StatusCode() != 200 { // 204 No Content is expected, but sometimes 200
		return fmt.Errorf("request failed: %s", resp.Status())
	}

	fmt.Printf("Revoked %s on %s from role %s\n", crGrantPrivilege, crGrantType, catalogRoleName)
	return nil
}

func runCatalogRolesGrantsList(cmd *cobra.Command, args []string) error {
	client, _, err := newManagementClient()
	if err != nil {
		return err
	}

	resp, err := client.ListGrantsForCatalogRoleWithResponse(context.Background(), catalogName, catalogRoleName)
	if err != nil {
		return err
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("request failed: %s", resp.Status())
	}

	// Parse raw body manually because generated types are incomplete
	var result struct {
		Grants []map[string]interface{} `json:"grants"`
	}
	if err := json.Unmarshal(resp.Body, &result); err != nil {
		return fmt.Errorf("failed to parse response: %w", err)
	}

	if len(result.Grants) == 0 {
		fmt.Println("(no grants)")
		return nil
	}

	for _, g := range result.Grants {
		typ, _ := g["type"].(string)
		priv, _ := g["privilege"].(string)

		fmt.Printf("Type: %s, Privilege: %s", typ, priv)

		if ns, ok := g["namespace"].([]interface{}); ok {
			nsStrs := make([]string, len(ns))
			for i, v := range ns {
				nsStrs[i] = fmt.Sprint(v)
			}
			fmt.Printf(", Namespace: %s", strings.Join(nsStrs, "."))
		}
		if tbl, ok := g["tableName"].(string); ok {
			fmt.Printf(", Table: %s", tbl)
		}
		if view, ok := g["viewName"].(string); ok {
			fmt.Printf(", View: %s", view)
		}
		fmt.Println()
	}

	return nil
}

func runCatalogRolesPrincipalsList(cmd *cobra.Command, args []string) error {
	client, _, err := newManagementClient()
	if err != nil {
		return err
	}

	resp, err := client.ListAssigneePrincipalRolesForCatalogRoleWithResponse(context.Background(), catalogName, catalogRoleName)
	if err != nil {
		return err
	}
	if resp.JSON200 == nil {
		return fmt.Errorf("request failed: %s", resp.Status())
	}

	if len(resp.JSON200.Roles) == 0 {
		fmt.Println("(no principal roles assigned)")
		return nil
	}

	for _, r := range resp.JSON200.Roles {
		fmt.Println(r.Name)
	}

	return nil
}

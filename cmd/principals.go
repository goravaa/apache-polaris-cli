package cmd

import (
	"context"
	"encoding/json"
	"fmt"

	managementapi "github.com/goravaa/apache-polaris-cli/pkg/api/openapi/management"
	"github.com/spf13/cobra"
)

var (
	principalName       string
	principalRole       string
	principalProperties []string
)

var principalsCmd = &cobra.Command{
	Use:   "principals",
	Short: "Principal management commands",
}

var principalsListCmd = &cobra.Command{
	Use:   "list",
	Short: "List principals",
	RunE:  runPrincipalsList,
}

var principalsCreateCmd = &cobra.Command{
	Use:   "create",
	Short: "Create a principal",
	RunE:  runPrincipalsCreate,
}

var principalsDeleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete a principal",
	RunE:  runPrincipalsDelete,
}

var principalsDescribeCmd = &cobra.Command{
	Use:   "describe",
	Short: "Describe a principal",
	RunE:  runPrincipalsDescribe,
}

var principalsUpdateCmd = &cobra.Command{
	Use:   "update",
	Short: "Update a principal",
	RunE:  runPrincipalsUpdate,
}

var principalsRotateCredentialsCmd = &cobra.Command{
	Use:   "rotate-credentials",
	Short: "Rotate principal credentials",
	RunE:  runPrincipalsRotateCredentials,
}

var principalsResetCredentialsCmd = &cobra.Command{
	Use:   "reset-credentials",
	Short: "Reset principal credentials",
	RunE:  runPrincipalsResetCredentials,
}

var principalsListRolesCmd = &cobra.Command{
	Use:   "list-roles",
	Short: "List roles assigned to a principal",
	RunE:  runPrincipalsListRoles,
}

var principalsAssignRoleCmd = &cobra.Command{
	Use:   "assign-role",
	Short: "Assign a role to a principal",
	RunE:  runPrincipalsAssignRole,
}

var principalsRevokeRoleCmd = &cobra.Command{
	Use:   "revoke-role",
	Short: "Revoke a role from a principal",
	RunE:  runPrincipalsRevokeRole,
}

func init() {
	rootCmd.AddCommand(principalsCmd)
	principalsCmd.AddCommand(principalsListCmd)
	principalsCmd.AddCommand(principalsCreateCmd)
	principalsCmd.AddCommand(principalsDeleteCmd)
	principalsCmd.AddCommand(principalsDescribeCmd)
	principalsCmd.AddCommand(principalsUpdateCmd)
	principalsCmd.AddCommand(principalsRotateCredentialsCmd)
	principalsCmd.AddCommand(principalsResetCredentialsCmd)
	principalsCmd.AddCommand(principalsListRolesCmd)
	principalsCmd.AddCommand(principalsAssignRoleCmd)
	principalsCmd.AddCommand(principalsRevokeRoleCmd)

	principalsCreateCmd.Flags().StringVar(&principalName, "name", "", "Principal name (required)")
	principalsCreateCmd.Flags().StringArrayVar(&principalProperties, "property", nil, "Principal property key=value (repeatable)")

	principalsDeleteCmd.Flags().StringVar(&principalName, "name", "", "Principal name (required)")

	principalsDescribeCmd.Flags().StringVar(&principalName, "name", "", "Principal name (required)")

	principalsUpdateCmd.Flags().StringVar(&principalName, "name", "", "Principal name (required)")
	principalsUpdateCmd.Flags().StringArrayVar(&principalProperties, "property", nil, "Principal property key=value (repeatable)")

	principalsRotateCredentialsCmd.Flags().StringVar(&principalName, "name", "", "Principal name (required)")

	principalsResetCredentialsCmd.Flags().StringVar(&principalName, "name", "", "Principal name (required)")

	principalsListRolesCmd.Flags().StringVar(&principalName, "name", "", "Principal name (required)")

	principalsAssignRoleCmd.Flags().StringVar(&principalName, "name", "", "Principal name (required)")
	principalsAssignRoleCmd.Flags().StringVar(&principalRole, "role", "", "Principal role name (required)")

	principalsRevokeRoleCmd.Flags().StringVar(&principalName, "name", "", "Principal name (required)")
	principalsRevokeRoleCmd.Flags().StringVar(&principalRole, "role", "", "Principal role name (required)")
}

func runPrincipalsList(cmd *cobra.Command, args []string) error {
	client, _, err := newManagementClient()
	if err != nil {
		return err
	}

	resp, err := client.ListPrincipalsWithResponse(context.Background())
	if err != nil {
		return err
	}
	if resp.JSON200 == nil {
		return fmt.Errorf("request failed: %s", resp.Status())
	}

	if len(resp.JSON200.Principals) == 0 {
		fmt.Println("(no principals)")
		return nil
	}

	for _, p := range resp.JSON200.Principals {
		fmt.Println(p.Name)
	}

	return nil
}

func runPrincipalsCreate(cmd *cobra.Command, args []string) error {
	if principalName == "" {
		return fmt.Errorf("--name is required")
	}

	props, err := parseProperties(principalProperties)
	if err != nil {
		return err
	}

	req := managementapi.CreatePrincipalRequest{
		Principal: &managementapi.Principal{
			Name: principalName,
		},
	}
	if len(props) > 0 {
		req.Principal.Properties = &props
	}

	client, _, err := newManagementClient()
	if err != nil {
		return err
	}

	resp, err := client.CreatePrincipalWithResponse(context.Background(), req)
	if err != nil {
		return err
	}
	if resp.JSON201 == nil {
		return fmt.Errorf("request failed: %s", resp.Status())
	}

	bytes, err := json.MarshalIndent(resp.JSON201, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}
	fmt.Println(string(bytes))

	return nil
}

func runPrincipalsDelete(cmd *cobra.Command, args []string) error {
	if principalName == "" {
		return fmt.Errorf("--name is required")
	}

	client, _, err := newManagementClient()
	if err != nil {
		return err
	}

	resp, err := client.DeletePrincipalWithResponse(context.Background(), principalName)
	if err != nil {
		return err
	}
	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		return fmt.Errorf("request failed: %s", resp.Status())
	}

	fmt.Printf("Deleted principal %s\n", principalName)
	return nil
}

func runPrincipalsDescribe(cmd *cobra.Command, args []string) error {
	if principalName == "" {
		return fmt.Errorf("--name is required")
	}

	client, _, err := newManagementClient()
	if err != nil {
		return err
	}

	resp, err := client.GetPrincipalWithResponse(context.Background(), principalName)
	if err != nil {
		return err
	}
	if resp.JSON200 == nil {
		return fmt.Errorf("request failed: %s", resp.Status())
	}

	p := resp.JSON200
	fmt.Printf("Name: %s\n", p.Name)
	if p.ClientId != nil {
		fmt.Printf("Client ID: %s\n", *p.ClientId)
	}
	if p.EntityVersion != nil {
		fmt.Printf("Entity Version: %d\n", *p.EntityVersion)
	}
	if p.CreateTimestamp != nil {
		fmt.Printf("Create Timestamp: %d\n", *p.CreateTimestamp)
	}
	if p.LastUpdateTimestamp != nil {
		fmt.Printf("Last Update Timestamp: %d\n", *p.LastUpdateTimestamp)
	}
	if p.Properties != nil && len(*p.Properties) > 0 {
		fmt.Println("Properties:")
		for k, v := range *p.Properties {
			fmt.Printf("  %s=%s\n", k, v)
		}
	}

	return nil
}

func runPrincipalsUpdate(cmd *cobra.Command, args []string) error {
	if principalName == "" {
		return fmt.Errorf("--name is required")
	}

	props, err := parseProperties(principalProperties)
	if err != nil {
		return err
	}

	client, _, err := newManagementClient()
	if err != nil {
		return err
	}

	// Fetch current principal to get entity version and existing properties
	getResp, err := client.GetPrincipalWithResponse(context.Background(), principalName)
	if err != nil {
		return err
	}
	if getResp.JSON200 == nil {
		return fmt.Errorf("failed to fetch principal: %s", getResp.Status())
	}

	current := getResp.JSON200
	if current.EntityVersion == nil {
		return fmt.Errorf("principal has no entity version")
	}

	// Merge properties
	updatedProps := make(map[string]string)
	if current.Properties != nil {
		for k, v := range *current.Properties {
			updatedProps[k] = v
		}
	}
	for k, v := range props {
		updatedProps[k] = v
	}

	req := managementapi.UpdatePrincipalRequest{
		CurrentEntityVersion: *current.EntityVersion,
		Properties:           updatedProps,
	}

	resp, err := client.UpdatePrincipalWithResponse(context.Background(), principalName, req)
	if err != nil {
		return err
	}
	if resp.JSON200 == nil {
		return fmt.Errorf("request failed: %s", resp.Status())
	}

	fmt.Printf("Updated principal %s\n", principalName)
	return nil
}

func runPrincipalsRotateCredentials(cmd *cobra.Command, args []string) error {
	if principalName == "" {
		return fmt.Errorf("--name is required")
	}

	client, _, err := newManagementClient()
	if err != nil {
		return err
	}

	resp, err := client.RotateCredentialsWithResponse(context.Background(), principalName)
	if err != nil {
		return err
	}
	if resp.JSON200 == nil {
		return fmt.Errorf("request failed: %s", resp.Status())
	}

	bytes, err := json.MarshalIndent(resp.JSON200, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}
	fmt.Println(string(bytes))

	return nil
}

func runPrincipalsResetCredentials(cmd *cobra.Command, args []string) error {
	if principalName == "" {
		return fmt.Errorf("--name is required")
	}

	client, _, err := newManagementClient()
	if err != nil {
		return err
	}

	// Empty body implies server generation
	req := managementapi.ResetPrincipalRequest{}

	resp, err := client.ResetCredentialsWithResponse(context.Background(), principalName, req)
	if err != nil {
		return err
	}
	if resp.JSON200 == nil {
		return fmt.Errorf("request failed: %s", resp.Status())
	}

	bytes, err := json.MarshalIndent(resp.JSON200, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal response: %w", err)
	}
	fmt.Println(string(bytes))

	return nil
}

func runPrincipalsListRoles(cmd *cobra.Command, args []string) error {
	if principalName == "" {
		return fmt.Errorf("--name is required")
	}

	client, _, err := newManagementClient()
	if err != nil {
		return err
	}

	resp, err := client.ListPrincipalRolesAssignedWithResponse(context.Background(), principalName)
	if err != nil {
		return err
	}
	if resp.JSON200 == nil {
		return fmt.Errorf("request failed: %s", resp.Status())
	}

	if len(resp.JSON200.Roles) == 0 {
		fmt.Println("(no roles assigned)")
		return nil
	}

	for _, r := range resp.JSON200.Roles {
		fmt.Println(r.Name)
	}

	return nil
}

func runPrincipalsAssignRole(cmd *cobra.Command, args []string) error {
	if principalName == "" {
		return fmt.Errorf("--name is required")
	}
	if principalRole == "" {
		return fmt.Errorf("--role is required")
	}

	client, _, err := newManagementClient()
	if err != nil {
		return err
	}

	req := managementapi.GrantPrincipalRoleRequest{
		PrincipalRole: &managementapi.PrincipalRole{
			Name: principalRole,
		},
	}

	resp, err := client.AssignPrincipalRoleWithResponse(context.Background(), principalName, req)
	if err != nil {
		return err
	}
	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		return fmt.Errorf("request failed: %s", resp.Status())
	}

	fmt.Printf("Assigned role %s to principal %s\n", principalRole, principalName)
	return nil
}

func runPrincipalsRevokeRole(cmd *cobra.Command, args []string) error {
	if principalName == "" {
		return fmt.Errorf("--name is required")
	}
	if principalRole == "" {
		return fmt.Errorf("--role is required")
	}

	client, _, err := newManagementClient()
	if err != nil {
		return err
	}

	resp, err := client.RevokePrincipalRoleWithResponse(context.Background(), principalName, principalRole)
	if err != nil {
		return err
	}
	if resp.StatusCode() < 200 || resp.StatusCode() >= 300 {
		return fmt.Errorf("request failed: %s", resp.Status())
	}

	fmt.Printf("Revoked role %s from principal %s\n", principalRole, principalName)
	return nil
}

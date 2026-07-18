//go:build integration
// +build integration

package integration

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var (
	binaryName = "polaris"
	binaryPath string
)

func TestMain(m *testing.M) {
	err := buildCLI()
	if err != nil {
		fmt.Printf("could not build CLI: %v\n", err)
		os.Exit(1)
	}

	code := m.Run()
	os.Remove(binaryPath)
	os.Exit(code)
}

func buildCLI() error {
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	// Go up to the root directory
	rootDir := filepath.Dir(filepath.Dir(dir))

	binaryPath = filepath.Join(rootDir, binaryName)
	cmd := exec.Command("go", "build", "-o", binaryPath, rootDir)
	return cmd.Run()
}

func runCLI(t *testing.T, tmpDir string, args ...string) (string, error) {
	cmd := exec.Command(binaryPath, args...)

	// Create a temporary directory for config to isolate tests
	cmd.Env = append(os.Environ(), fmt.Sprintf("HOME=%s", tmpDir))

	out, err := cmd.CombinedOutput()
	return string(out), err
}

func TestCLIIntegration(t *testing.T) {
	host := os.Getenv("POLARIS_TEST_HOST")
	clientID := os.Getenv("POLARIS_TEST_CLIENT_ID")
	clientSecret := os.Getenv("POLARIS_TEST_CLIENT_SECRET")
	require.NotEmpty(t, host, "POLARIS_TEST_HOST must be set")

	// Setup isolated home dir
	tmpDir := t.TempDir()

	// Setup config
	_, err := runCLI(t, tmpDir, "config", "set", "--host", host)
	require.NoError(t, err)

	t.Run("Test Auth Login", func(t *testing.T) {
		out, err := runCLI(t, tmpDir, "auth", "login", "--client-id", clientID, "--client-secret", clientSecret)
		require.NoError(t, err, "Auth login failed: %s", out)
		assert.Contains(t, out, "Successfully authenticated")
	})

	t.Run("Test Principal Lifecycle", func(t *testing.T) {
		principalName := fmt.Sprintf("test-principal-%d", time.Now().UnixNano())

		// Create
		out, err := runCLI(t, tmpDir, "principals", "create", "--name", principalName)
		require.NoError(t, err, "Principal creation failed: %s", out)
		assert.Contains(t, out, principalName)

		// List
		out, err = runCLI(t, tmpDir, "principals", "list")
		require.NoError(t, err, "Principal list failed: %s", out)
		assert.Contains(t, out, principalName)

		// Get
		out, err = runCLI(t, tmpDir, "principals", "get", principalName)
		require.NoError(t, err, "Principal get failed: %s", out)
		assert.Contains(t, out, principalName)

		// Delete
		out, err = runCLI(t, tmpDir, "principals", "delete", principalName)
		require.NoError(t, err, "Principal deletion failed: %s", out)
		assert.Contains(t, out, "Successfully deleted")

		// Verify deletion
		out, err = runCLI(t, tmpDir, "principals", "get", principalName)
		require.Error(t, err, "Principal should be deleted")
	})

	t.Run("Test Catalog Lifecycle", func(t *testing.T) {
		catalogName := fmt.Sprintf("test-catalog-%d", time.Now().UnixNano())

		// Create
		out, err := runCLI(t, tmpDir, "catalogs", "create", "--name", catalogName, "--type", "INTERNAL", "--default-base-location", "s3://test-bucket/")
		require.NoError(t, err, "Catalog creation failed: %s", out)
		assert.Contains(t, out, catalogName)

		// List
		out, err = runCLI(t, tmpDir, "catalogs", "list")
		require.NoError(t, err, "Catalog list failed: %s", out)
		assert.Contains(t, out, catalogName)

		// Get
		out, err = runCLI(t, tmpDir, "catalogs", "get", catalogName)
		require.NoError(t, err, "Catalog get failed: %s", out)
		assert.Contains(t, out, catalogName)

		// Delete
		out, err = runCLI(t, tmpDir, "catalogs", "delete", catalogName)
		require.NoError(t, err, "Catalog deletion failed: %s", out)
		assert.Contains(t, out, "Successfully deleted")

		// Verify deletion
		out, err = runCLI(t, tmpDir, "catalogs", "get", catalogName)
		require.Error(t, err, "Catalog should be deleted")
	})

	t.Run("Test Advanced Auth and Config", func(t *testing.T) {
		// Test config show
		out, err := runCLI(t, tmpDir, "config", "show")
		require.NoError(t, err)
		assert.Contains(t, out, host)

		// Test auth status
		out, err = runCLI(t, tmpDir, "auth", "status")
		require.NoError(t, err)
		assert.Contains(t, out, "Currently authenticated")

		// Test auth refresh
		out, err = runCLI(t, tmpDir, "auth", "refresh")
		require.NoError(t, err)
		assert.Contains(t, out, "Successfully refreshed")

		// We do not test logout yet because it breaks subsequent tests unless we login again.
	})

	t.Run("Test Principal Advanced Lifecycle", func(t *testing.T) {
		principalName := fmt.Sprintf("test-adv-principal-%d", time.Now().UnixNano())
		out, err := runCLI(t, tmpDir, "principals", "create", "--name", principalName)
		require.NoError(t, err)

		// Describe
		out, err = runCLI(t, tmpDir, "principals", "describe", principalName)
		require.NoError(t, err)
		assert.Contains(t, out, principalName)

		// Update (e.g. changing property)
		out, err = runCLI(t, tmpDir, "principals", "update", principalName, "--property", "newProp=newValue")
		require.NoError(t, err)
		assert.Contains(t, out, "Successfully updated")

		// Rotate Credentials
		out, err = runCLI(t, tmpDir, "principals", "rotate-credentials", principalName)
		require.NoError(t, err)
		assert.Contains(t, out, "Client ID")
		assert.Contains(t, out, "Client Secret")

		// Reset Credentials
		out, err = runCLI(t, tmpDir, "principals", "reset-credentials", principalName)
		require.NoError(t, err)
		assert.Contains(t, out, "Successfully reset")

		// Roles (assuming empty but checking command works)
		out, err = runCLI(t, tmpDir, "principals", "list-roles", principalName)
		require.NoError(t, err)
		// Should not error, output might just be empty list or headers

		// Delete
		_, err = runCLI(t, tmpDir, "principals", "delete", principalName)
		require.NoError(t, err)
	})

	t.Run("Test Catalog, Namespaces and Tables", func(t *testing.T) {
		catalogName := fmt.Sprintf("test-ns-catalog-%d", time.Now().UnixNano())

		// 1. Create Catalog
		_, err := runCLI(t, tmpDir, "catalogs", "create", "--name", catalogName, "--type", "INTERNAL", "--default-base-location", "s3://test-bucket/")
		require.NoError(t, err)

		// Describe
		out, err := runCLI(t, tmpDir, "catalogs", "describe", catalogName)
		require.NoError(t, err)
		assert.Contains(t, out, catalogName)

		// 2. Namespaces
		namespaceName := "test_namespace"

		// Create Namespace
		out, err = runCLI(t, tmpDir, "catalog", "namespaces", "create", namespaceName, "--catalog", catalogName)
		require.NoError(t, err)
		assert.Contains(t, out, namespaceName)

		// List Namespaces
		out, err = runCLI(t, tmpDir, "catalog", "namespaces", "list", "--catalog", catalogName)
		require.NoError(t, err)
		assert.Contains(t, out, namespaceName)

		// 3. Tables (List only, as create is not yet implemented)
		out, err = runCLI(t, tmpDir, "catalog", "tables", "list", "--catalog", catalogName, "--namespace", namespaceName)
		require.NoError(t, err) // Should succeed, returning empty list

		// 4. Delete Namespace
		out, err = runCLI(t, tmpDir, "catalog", "namespaces", "delete", namespaceName, "--catalog", catalogName)
		require.NoError(t, err)
		assert.Contains(t, out, "Successfully deleted")

		// 5. Delete Catalog
		_, err = runCLI(t, tmpDir, "catalogs", "delete", catalogName)
		require.NoError(t, err)
	})

	t.Run("Test Auth Failure", func(t *testing.T) {
		tmpDirFail := t.TempDir()
		_, err := runCLI(t, tmpDirFail, "config", "set", "--host", host)
		require.NoError(t, err)

		out, err := runCLI(t, tmpDirFail, "auth", "login", "--client-id", "invalid", "--client-secret", "invalid")
		require.Error(t, err, "Expected failure with invalid credentials")
		assert.Contains(t, string(out), "invalid_client")
	})
}

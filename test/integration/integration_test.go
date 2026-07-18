//go:build integration
// +build integration

package integration

import (
	"encoding/json"
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

		// Describe
		out, err = runCLI(t, tmpDir, "principals", "describe", "--name", principalName)
		require.NoError(t, err, "Principal describe failed: %s", out)
		assert.Contains(t, out, principalName)

		// Delete
		out, err = runCLI(t, tmpDir, "principals", "delete", "--name", principalName)
		require.NoError(t, err, "Principal deletion failed: %s", out)
		assert.Contains(t, out, "Deleted principal")

		// Verify deletion
		out, err = runCLI(t, tmpDir, "principals", "describe", "--name", principalName)
		require.Error(t, err, "Principal should be deleted")
	})

	t.Run("Test Catalog Lifecycle", func(t *testing.T) {
		catalogName := fmt.Sprintf("test-catalog-%d", time.Now().UnixNano())
		baseLocation := fmt.Sprintf("s3://test-bucket/%s/", catalogName)

		// Create
		out, err := runCLI(t, tmpDir, "catalogs", "create", "--name", catalogName, "--type", "INTERNAL", "--default-base-location", baseLocation)
		require.NoError(t, err, "Catalog creation failed: %s", out)
		assert.Contains(t, out, catalogName)

		// List
		out, err = runCLI(t, tmpDir, "catalogs", "list")
		require.NoError(t, err, "Catalog list failed: %s", out)
		assert.Contains(t, out, catalogName)

		// Describe
		out, err = runCLI(t, tmpDir, "catalogs", "describe", "--name", catalogName)
		require.NoError(t, err, "Catalog describe failed: %s", out)
		assert.Contains(t, out, catalogName)

		// Delete
		out, err = runCLI(t, tmpDir, "catalogs", "delete", "--name", catalogName)
		require.NoError(t, err, "Catalog deletion failed: %s", out)
		assert.Contains(t, out, "Deleted catalog")

		// Verify deletion
		out, err = runCLI(t, tmpDir, "catalogs", "describe", "--name", catalogName)
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
		assert.Contains(t, out, "Status: Authenticated")

		// Test auth refresh
		out, err = runCLI(t, tmpDir, "auth", "refresh")
		require.NoError(t, err)
		assert.Contains(t, out, "Token refreshed successfully")

		// We do not test logout yet because it breaks subsequent tests unless we login again.
	})

	t.Run("Test Principal Advanced Lifecycle", func(t *testing.T) {
		principalName := fmt.Sprintf("test-adv-principal-%d", time.Now().UnixNano())
		out, err := runCLI(t, tmpDir, "principals", "create", "--name", principalName)
		require.NoError(t, err)

		var created struct {
			Credentials struct {
				ClientID     string `json:"clientId"`
				ClientSecret string `json:"clientSecret"`
			} `json:"credentials"`
		}
		require.NoError(t, json.Unmarshal([]byte(out), &created))
		require.NotEmpty(t, created.Credentials.ClientID)
		require.NotEmpty(t, created.Credentials.ClientSecret)

		// Describe
		out, err = runCLI(t, tmpDir, "principals", "describe", "--name", principalName)
		require.NoError(t, err)
		assert.Contains(t, out, principalName)

		// Update (e.g. changing property)
		out, err = runCLI(t, tmpDir, "principals", "update", "--name", principalName, "--property", "newProp=newValue")
		require.NoError(t, err)
		assert.Contains(t, out, "Updated principal")

		// Rotate Credentials must be done as the principal itself (root gets 403).
		principalHome := t.TempDir()
		_, err = runCLI(t, principalHome, "config", "set", "--host", host)
		require.NoError(t, err)
		_, err = runCLI(t, principalHome, "auth", "login",
			"--client-id", created.Credentials.ClientID,
			"--client-secret", created.Credentials.ClientSecret,
		)
		require.NoError(t, err)
		out, err = runCLI(t, principalHome, "principals", "rotate-credentials", "--name", principalName)
		require.NoError(t, err)
		assert.Contains(t, out, "clientId")
		assert.Contains(t, out, "clientSecret")

		// Reset Credentials (admin operation)
		out, err = runCLI(t, tmpDir, "principals", "reset-credentials", "--name", principalName)
		require.NoError(t, err)
		assert.Contains(t, out, "clientId")
		assert.Contains(t, out, "clientSecret")

		// Roles (assuming empty but checking command works)
		out, err = runCLI(t, tmpDir, "principals", "list-roles", "--name", principalName)
		require.NoError(t, err)
		// Should not error, output might just be empty list or headers

		// Delete
		_, err = runCLI(t, tmpDir, "principals", "delete", "--name", principalName)
		require.NoError(t, err)
	})

	t.Run("Test Catalog, Namespaces and Tables", func(t *testing.T) {
		catalogName := fmt.Sprintf("test-ns-catalog-%d", time.Now().UnixNano())
		baseLocation := fmt.Sprintf("s3://test-bucket/%s/", catalogName)

		// 1. Create Catalog
		_, err := runCLI(t, tmpDir, "catalogs", "create", "--name", catalogName, "--type", "INTERNAL", "--default-base-location", baseLocation)
		require.NoError(t, err)

		// Describe
		out, err := runCLI(t, tmpDir, "catalogs", "describe", "--name", catalogName)
		require.NoError(t, err)
		assert.Contains(t, out, catalogName)

		// 2. Namespaces
		namespaceName := "test_namespace"

		// Create Namespace
		out, err = runCLI(t, tmpDir, "catalog", "namespaces", "create", namespaceName, "--prefix", catalogName)
		require.NoError(t, err)
		assert.Contains(t, out, namespaceName)

		// List Namespaces
		out, err = runCLI(t, tmpDir, "catalog", "namespaces", "list", "--prefix", catalogName)
		require.NoError(t, err)
		assert.Contains(t, out, namespaceName)

		// 3. Tables (List only, as create is not yet implemented)
		out, err = runCLI(t, tmpDir, "catalog", "tables", "list", "--prefix", catalogName, "--namespace", namespaceName)
		require.NoError(t, err) // Should succeed, returning empty list

		// 4. Delete Namespace
		out, err = runCLI(t, tmpDir, "catalog", "namespaces", "delete", namespaceName, "--prefix", catalogName)
		require.NoError(t, err)
		assert.Contains(t, out, "Deleted namespace")

		// 5. Delete Catalog
		_, err = runCLI(t, tmpDir, "catalogs", "delete", "--name", catalogName)
		require.NoError(t, err)
	})

	t.Run("Test Auth Failure", func(t *testing.T) {
		tmpDirFail := t.TempDir()
		_, err := runCLI(t, tmpDirFail, "config", "set", "--host", host)
		require.NoError(t, err)

		out, err := runCLI(t, tmpDirFail, "auth", "login", "--client-id", "invalid", "--client-secret", "invalid")
		require.Error(t, err, "Expected failure with invalid credentials")
		assert.Contains(t, string(out), "authentication failed")
	})
}

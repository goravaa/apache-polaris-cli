package cmd

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/goravaa/apache-polaris-cli/pkg/api"
	catalogapi "github.com/goravaa/apache-polaris-cli/pkg/api/openapi/catalog"
	"github.com/goravaa/apache-polaris-cli/pkg/config"
	"github.com/spf13/cobra"
)

var catalogPrefix string

var catalogCmd = &cobra.Command{
	Use:   "catalog",
	Short: "Catalog API commands",
}

func init() {
	rootCmd.AddCommand(catalogCmd)
	catalogCmd.PersistentFlags().StringVar(&catalogPrefix, "prefix", "", "Catalog prefix (required if not set in config)")
}

func newCatalogClient() (*catalogapi.ClientWithResponses, *config.Config, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load config: %w", err)
	}

	token, err := config.GetAccessToken()
	if err != nil {
		return nil, nil, err
	}

	baseURL := strings.TrimSuffix(cfg.Host, "/") + "/api/catalog"
	httpClient := &http.Client{Timeout: 30 * time.Second}

	client, err := catalogapi.NewClientWithResponses(
		baseURL,
		catalogapi.WithHTTPClient(httpClient),
		catalogapi.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			if cfg.Realm != "" {
				req.Header.Set(api.RealmHeaderName, cfg.Realm)
			}
			return nil
		}),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create catalog client: %w", err)
	}

	return client, cfg, nil
}

func resolveCatalogPrefix(cfg *config.Config) (string, error) {
	if catalogPrefix != "" {
		return catalogPrefix, nil
	}
	if cfg.CatalogPrefix != "" {
		return cfg.CatalogPrefix, nil
	}
	return "", fmt.Errorf("catalog prefix is required. Use --prefix or set --catalog-prefix in config")
}

func parseNamespaceArg(input string) ([]string, error) {
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return nil, fmt.Errorf("namespace is required")
	}
	if strings.Contains(trimmed, "\x1f") {
		return strings.Split(trimmed, "\x1f"), nil
	}
	if strings.Contains(trimmed, "%1F") {
		decoded := strings.ReplaceAll(trimmed, "%1F", "\x1f")
		return strings.Split(decoded, "\x1f"), nil
	}
	if strings.Contains(trimmed, "/") {
		return strings.Split(trimmed, "/"), nil
	}
	return strings.Split(trimmed, "."), nil
}

func namespacePath(parts []string) string {
	return strings.Join(parts, "\x1f")
}

func formatNamespace(parts []string) string {
	return strings.Join(parts, ".")
}

func parseProperties(pairs []string) (map[string]string, error) {
	if len(pairs) == 0 {
		return nil, nil
	}
	props := make(map[string]string, len(pairs))
	for _, pair := range pairs {
		kv := strings.SplitN(pair, "=", 2)
		if len(kv) != 2 || strings.TrimSpace(kv[0]) == "" {
			return nil, fmt.Errorf("invalid property %q (expected key=value)", pair)
		}
		props[strings.TrimSpace(kv[0])] = kv[1]
	}
	return props, nil
}

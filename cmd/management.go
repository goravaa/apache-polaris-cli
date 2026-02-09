package cmd

import (
	"context"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/goravaa/apache-polaris-cli/pkg/api"
	managementapi "github.com/goravaa/apache-polaris-cli/pkg/api/openapi/management"
	"github.com/goravaa/apache-polaris-cli/pkg/config"
)

func newManagementClient() (*managementapi.ClientWithResponses, *config.Config, error) {
	cfg, err := config.LoadConfig()
	if err != nil {
		return nil, nil, fmt.Errorf("failed to load config: %w", err)
	}

	token, err := config.GetAccessToken()
	if err != nil {
		return nil, nil, err
	}

	baseURL := strings.TrimSuffix(cfg.Host, "/") + "/api/management/v1"
	httpClient := &http.Client{Timeout: 30 * time.Second}

	client, err := managementapi.NewClientWithResponses(
		baseURL,
		managementapi.WithHTTPClient(httpClient),
		managementapi.WithRequestEditorFn(func(ctx context.Context, req *http.Request) error {
			req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
			if cfg.Realm != "" {
				req.Header.Set(api.RealmHeaderName, cfg.Realm)
			}
			return nil
		}),
	)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create management client: %w", err)
	}

	return client, cfg, nil
}

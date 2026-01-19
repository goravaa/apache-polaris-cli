package api

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/goravaa/apache-polaris-cli/pkg/config"
)

const (
	TokenEndpoint   = "/api/catalog/v1/oauth/tokens"
	RealmHeaderName = "Polaris-Realm"
	DefaultScope    = "PRINCIPAL_ROLE:ALL"
)

type OAuthTokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType   string `json:"token_type"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
}

type OAuthError struct {
	Error            string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

type AuthClient struct {
	httpClient *http.Client
	config     *config.Config
}

func NewAuthClient(cfg *config.Config) *AuthClient {
	return &AuthClient{
		httpClient: &http.Client{
			Timeout: 30 * time.Second,
		},
		config: cfg,
	}
}

func (c *AuthClient) Login(clientID, clientSecret string) (*config.Credentials, error) {
	tokenURL := fmt.Sprintf("%s%s", strings.TrimSuffix(c.config.Host, "/"), TokenEndpoint)

	formData := url.Values{}
	formData.Set("grant_type", "client_credentials")
	formData.Set("client_id", clientID)
	formData.Set("client_secret", clientSecret)
	formData.Set("scope", DefaultScope)

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if c.config.Realm != "" {
		req.Header.Set(RealmHeaderName, c.config.Realm)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Polaris server: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var oauthErr OAuthError
		if err := json.Unmarshal(body, &oauthErr); err == nil && oauthErr.Error != "" {
			return nil, fmt.Errorf("authentication failed: %s - %s", oauthErr.Error, oauthErr.ErrorDescription)
		}
		return nil, fmt.Errorf("authentication failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp OAuthTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	if tokenResp.AccessToken == "" {
		return nil, fmt.Errorf("received empty access token from server")
	}

	credentials := &config.Credentials{
		AccessToken:  tokenResp.AccessToken,
		TokenType:    tokenResp.TokenType,
		ExpiresIn:    tokenResp.ExpiresIn,
		Scope:        tokenResp.Scope,
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}

	return credentials, nil
}

func (c *AuthClient) RefreshToken(currentToken string) (*config.Credentials, error) {
	tokenURL := fmt.Sprintf("%s%s", strings.TrimSuffix(c.config.Host, "/"), TokenEndpoint)

	formData := url.Values{}
	formData.Set("grant_type", "urn:ietf:params:oauth:grant-type:token-exchange")
	formData.Set("subject_token", currentToken)
	formData.Set("subject_token_type", "urn:ietf:params:oauth:token-type:access_token")

	req, err := http.NewRequest("POST", tokenURL, strings.NewReader(formData.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	if c.config.Realm != "" {
		req.Header.Set(RealmHeaderName, c.config.Realm)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to Polaris server: %w", err)
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response: %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		var oauthErr OAuthError
		if err := json.Unmarshal(body, &oauthErr); err == nil && oauthErr.Error != "" {
			return nil, fmt.Errorf("token refresh failed: %s - %s", oauthErr.Error, oauthErr.ErrorDescription)
		}
		return nil, fmt.Errorf("token refresh failed with status %d: %s", resp.StatusCode, string(body))
	}

	var tokenResp OAuthTokenResponse
	if err := json.Unmarshal(body, &tokenResp); err != nil {
		return nil, fmt.Errorf("failed to parse token response: %w", err)
	}

	if tokenResp.AccessToken == "" {
		return nil, fmt.Errorf("received empty access token from server")
	}

	existingCreds, _ := config.LoadCredentials()

	credentials := &config.Credentials{
		AccessToken: tokenResp.AccessToken,
		TokenType:   tokenResp.TokenType,
		ExpiresIn:   tokenResp.ExpiresIn,
		Scope:       tokenResp.Scope,
	}

	if existingCreds != nil {
		credentials.ClientID = existingCreds.ClientID
		credentials.ClientSecret = existingCreds.ClientSecret
	}

	return credentials, nil
}

func (c *AuthClient) Logout() error {
	return config.ClearCredentials()
}

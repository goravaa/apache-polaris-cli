package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

const (
	ConfigDirName       = ".polaris-cli"
	ConfigFileName      = "config.json"
	CredentialsFileName = "credentials.json"
)

type Config struct {
	Host  string `json:"host"`
	Realm string `json:"realm"`
}

type Credentials struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

func getConfigDir() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}
	return filepath.Join(homeDir, ConfigDirName), nil
}

func ensureConfigDir() (string, error) {
	configDir, err := getConfigDir()
	if err != nil {
		return "", err
	}

	if err := os.MkdirAll(configDir, 0700); err != nil {
		return "", fmt.Errorf("failed to create config directory: %w", err)
	}

	return configDir, nil
}

func LoadConfig() (*Config, error) {
	configDir, err := getConfigDir()
	if err != nil {
		return nil, err
	}

	configPath := filepath.Join(configDir, ConfigFileName)
	data, err := os.ReadFile(configPath)
	if err != nil {
		if os.IsNotExist(err) {
			return &Config{
				Host: "http://localhost:8181",
			}, nil
		}
		return nil, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config
	if err := json.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to parse config file: %w", err)
	}

	return &config, nil
}

func SaveConfig(config *Config) error {
	configDir, err := ensureConfigDir()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(config, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal config: %w", err)
	}

	configPath := filepath.Join(configDir, ConfigFileName)
	if err := os.WriteFile(configPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write config file: %w", err)
	}

	return nil
}

func LoadCredentials() (*Credentials, error) {
	configDir, err := getConfigDir()
	if err != nil {
		return nil, err
	}

	credentialsPath := filepath.Join(configDir, CredentialsFileName)
	data, err := os.ReadFile(credentialsPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("not authenticated. Please run 'polaris auth login' first")
		}
		return nil, fmt.Errorf("failed to read credentials file: %w", err)
	}

	var credentials Credentials
	if err := json.Unmarshal(data, &credentials); err != nil {
		return nil, fmt.Errorf("failed to parse credentials file: %w", err)
	}

	return &credentials, nil
}

func SaveCredentials(credentials *Credentials) error {
	configDir, err := ensureConfigDir()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(credentials, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal credentials: %w", err)
	}

	credentialsPath := filepath.Join(configDir, CredentialsFileName)
	if err := os.WriteFile(credentialsPath, data, 0600); err != nil {
		return fmt.Errorf("failed to write credentials file: %w", err)
	}

	return nil
}

func ClearCredentials() error {
	configDir, err := getConfigDir()
	if err != nil {
		return err
	}

	credentialsPath := filepath.Join(configDir, CredentialsFileName)
	if err := os.Remove(credentialsPath); err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("failed to remove credentials file: %w", err)
	}

	return nil
}

func IsAuthenticated() bool {
	creds, err := LoadCredentials()
	if err != nil {
		return false
	}
	return creds.AccessToken != ""
}

func GetAccessToken() (string, error) {
	creds, err := LoadCredentials()
	if err != nil {
		return "", err
	}
	if creds.AccessToken == "" {
		return "", fmt.Errorf("no access token found. Please run 'polaris auth login' first")
	}
	return creds.AccessToken, nil
}

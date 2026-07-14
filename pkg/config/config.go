package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

const (
	ConfigDirName       = ".polaris-cli"
	ConfigFileName      = "config.json"
	CredentialsFileName = "credentials.json"
)

type Config struct {
	Host             string `json:"host"`
	Realm            string `json:"realm"`
	CatalogPrefix    string `json:"catalog_prefix"`
	RootClientID     string `json:"root_client_id,omitempty"`
	RootClientSecret string `json:"root_client_secret,omitempty"`
}

func ApplyConfigValues(config *Config, values map[string]string) error {
	for key, value := range values {
		switch normalizeConfigKey(key) {
		case "host", "url", "server", "address", "serverurl", "serveraddress", "polarisserver", "polarisserverurl", "polarisserveraddress":
			config.Host = value
		case "realm", "polarisrealm":
			config.Realm = value
		case "catalogprefix", "prefix":
			config.CatalogPrefix = value
		case "rootclientid", "clientid":
			config.RootClientID = value
		case "rootclientsecret", "clientsecret", "secret":
			config.RootClientSecret = value
		}
	}
	return nil
}

func ParseConfigValues(data []byte) (map[string]string, error) {
	values, err := parseJSONConfigValues(data)
	if err == nil {
		return values, nil
	}

	values, yamlErr := parseYAMLConfigValues(data)
	if yamlErr == nil {
		return values, nil
	}

	return nil, fmt.Errorf("failed to parse config as JSON (%v) or YAML (%v)", err, yamlErr)
}

func parseJSONConfigValues(data []byte) (map[string]string, error) {
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, err
	}
	values := map[string]string{}
	flattenConfigValues("", raw, values)
	return values, nil
}

func flattenConfigValues(prefix string, raw map[string]any, values map[string]string) {
	for key, value := range raw {
		nextKey := key
		if prefix != "" {
			nextKey = prefix + "_" + key
		}

		switch typed := value.(type) {
		case map[string]any:
			flattenConfigValues(nextKey, typed, values)
		case string:
			values[nextKey] = typed
		case nil:
			values[nextKey] = ""
		default:
			values[nextKey] = fmt.Sprint(typed)
		}
	}
}

func parseYAMLConfigValues(data []byte) (map[string]string, error) {
	values := map[string]string{}
	lines := strings.Split(string(data), "\n")
	for lineNumber, line := range lines {
		line = stripYAMLComment(line)
		if strings.TrimSpace(line) == "" || strings.TrimSpace(line) == "---" {
			continue
		}
		if strings.HasPrefix(line, " ") || strings.HasPrefix(line, "\t") {
			return nil, fmt.Errorf("line %d: nested YAML is not supported", lineNumber+1)
		}

		key, value, ok := strings.Cut(line, ":")
		if !ok {
			return nil, fmt.Errorf("line %d: expected key: value", lineNumber+1)
		}

		key = strings.TrimSpace(key)
		if key == "" {
			return nil, fmt.Errorf("line %d: key is required", lineNumber+1)
		}
		values[key] = unquoteYAMLValue(strings.TrimSpace(value))
	}

	if len(values) == 0 {
		return nil, fmt.Errorf("no config values found")
	}
	return values, nil
}

func stripYAMLComment(line string) string {
	inSingleQuote := false
	inDoubleQuote := false

	for i, r := range line {
		switch r {
		case '\'':
			if !inDoubleQuote {
				inSingleQuote = !inSingleQuote
			}
		case '"':
			if !inSingleQuote {
				inDoubleQuote = !inDoubleQuote
			}
		case '#':
			if !inSingleQuote && !inDoubleQuote {
				return line[:i]
			}
		}
	}
	return line
}

func unquoteYAMLValue(value string) string {
	value = strings.TrimSpace(value)
	if len(value) < 2 {
		return value
	}
	if (value[0] == '"' && value[len(value)-1] == '"') || (value[0] == '\'' && value[len(value)-1] == '\'') {
		return value[1 : len(value)-1]
	}
	return value
}

func normalizeConfigKey(key string) string {
	replacer := strings.NewReplacer("_", "", "-", "", " ", "", ".", "")
	return strings.ToLower(replacer.Replace(key))
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

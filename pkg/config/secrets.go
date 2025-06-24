package config

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
)

// SecretsConfig represents the structure of the secrets configuration file
type SecretsConfig struct {
	Supertone struct {
		APIKey string `json:"api_key"`
		APIURL string `json:"api_url"`
	} `json:"supertone"`
	OtherProvider struct {
		APIKey string `json:"api_key"`
		APIURL string `json:"api_url"`
	} `json:"other_provider"`
}

// LoadSecrets loads the secrets configuration from the JSON file
func LoadSecrets() (*SecretsConfig, error) {
	// Get the current working directory
	wd, err := os.Getwd()
	if err != nil {
		return nil, fmt.Errorf("failed to get working directory: %w", err)
	}

	// Construct the path to the secrets file
	secretsPath := filepath.Join(wd, "config", "secrets", "api_keys.json")

	// Read the secrets file
	data, err := os.ReadFile(secretsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read secrets file %s: %w", secretsPath, err)
	}

	// Parse the JSON
	var secrets SecretsConfig
	if err := json.Unmarshal(data, &secrets); err != nil {
		return nil, fmt.Errorf("failed to parse secrets JSON: %w", err)
	}

	return &secrets, nil
}

// GetSupertoneAPIKey returns the Supertone API key from secrets
func GetSupertoneAPIKey() (string, error) {
	secrets, err := LoadSecrets()
	if err != nil {
		return "", err
	}
	return secrets.Supertone.APIKey, nil
}

// GetSupertoneAPIURL returns the Supertone API URL from secrets
func GetSupertoneAPIURL() (string, error) {
	secrets, err := LoadSecrets()
	if err != nil {
		return "", err
	}
	return secrets.Supertone.APIURL, nil
}

// GetProviderConfig returns the configuration for a specific provider
func GetProviderConfig(provider string) (string, string, error) {
	secrets, err := LoadSecrets()
	if err != nil {
		return "", "", err
	}

	switch provider {
	case "supertone":
		return secrets.Supertone.APIKey, secrets.Supertone.APIURL, nil
	case "other_provider":
		return secrets.OtherProvider.APIKey, secrets.OtherProvider.APIURL, nil
	default:
		return "", "", fmt.Errorf("unknown provider: %s", provider)
	}
} 
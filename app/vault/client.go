package vault

import (
	"fmt"
	"os"

	vault "github.com/hashicorp/vault/api"
	"poc-vault-go-kube/config"
)

type Client struct {
	*vault.Client
}

// NewClient creates a new Vault client using Kubernetes auth
func NewClient(cfg *config.Config) (*Client, error) {
	config := vault.DefaultConfig()
	config.Address = cfg.VaultAddr

	client, err := vault.NewClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create Vault client: %w", err)
	}

	// Read the service account token
	token, err := os.ReadFile("/var/run/secrets/kubernetes.io/serviceaccount/token")
	if err != nil {
		return nil, fmt.Errorf("failed to read service account token: %w", err)
	}

	// Login to Vault using Kubernetes auth
	_, err = client.Auth().Token().LookupSelf()
	if err != nil {
		// Token is not valid, try to authenticate
		resp, err := client.Logical().Write("auth/kubernetes/login", map[string]interface{}{
			"role": cfg.VaultRoleName,
			"jwt":  string(token),
		})
		if err != nil {
			return nil, fmt.Errorf("failed to authenticate with Vault: %w", err)
		}
		if resp == nil || resp.Auth == nil {
			return nil, fmt.Errorf("no auth info in response")
		}
		client.SetToken(resp.Auth.ClientToken)
	}

	return &Client{client}, nil
}

// GetClusterSecret retrieves a secret from the cluster-secrets path
func (c *Client) GetClusterSecret(path string) (map[string]interface{}, error) {
	secret, err := c.Logical().Read(fmt.Sprintf("cluster-secrets/data/%s", path))
	if err != nil {
		return nil, fmt.Errorf("failed to read cluster secret: %w", err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("secret not found")
	}
	
	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid secret format")
	}
	return data, nil
}

// GetAppSecret retrieves a secret from the app-secrets path
func (c *Client) GetAppSecret(path string) (map[string]interface{}, error) {
	secret, err := c.Logical().Read(fmt.Sprintf("app-secrets/data/%s", path))
	if err != nil {
		return nil, fmt.Errorf("failed to read app secret: %w", err)
	}
	if secret == nil || secret.Data == nil {
		return nil, fmt.Errorf("secret not found")
	}
	
	data, ok := secret.Data["data"].(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("invalid secret format")
	}
	return data, nil
}

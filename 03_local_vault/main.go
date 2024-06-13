package main

import (
	"fmt"
	"log"
	"os"
	"path"

	"github.com/hashicorp/vault/api"
	"gopkg.in/yaml.v3"
)

type Config struct {
	Vault struct {
		Address     string `yaml:"address"`
		RoleID      string `yaml:"role-id"`
		SecretID    string `yaml:"secret-id"`
		SecretsRoot string `yaml:"secrets-root"`
	} `yaml:"vault"`
}

func main() {
	// Read configuration from config.yaml
	configFile, err := os.ReadFile("config.yaml")
	if err != nil {
		log.Fatalf("Unable to read config file: %v", err)
	}

	var config Config
	err = yaml.Unmarshal(configFile, &config)
	if err != nil {
		log.Fatalf("Unable to parse config file: %v", err)
	}

	vaultAddr := config.Vault.Address
	roleID := config.Vault.RoleID
	secretID := config.Vault.SecretID
	secretPath := config.Vault.SecretsRoot

	// Initialize a Vault client
	vaultConfig := &api.Config{
		Address: vaultAddr,
	}
	client, err := api.NewClient(vaultConfig)
	if err != nil {
		log.Fatalf("Unable to initialize Vault client: %v", err)
	}

	// Auth with AppRole
	authData := map[string]interface{}{
		"role_id":   roleID,
		"secret_id": secretID,
	}
	authResp, err := client.Logical().Write("auth/approle/login", authData)
	if err != nil {
		log.Fatalf("Unable to authenticate with AppRole: %v", err)
	}
	client.SetToken(authResp.Auth.ClientToken)

	// Write the secret foo=bar
	key := "fizz"
	value := "buzz"
	secretData := map[string]interface{}{
		key: value,
	}
	fullPath := path.Join(secretPath, key)
	_, err = client.Logical().Write(fullPath, secretData)
	if err != nil {
		log.Fatalf("Unable to write secret: %v", err)
	}

	fmt.Println("Secret foo=bar written successfully")

	// Read the secret
	secret, err := client.Logical().Read(fullPath)
	if err != nil {
		log.Fatalf("Unable to read secret: %v", err)
	}

	if secret == nil {
		log.Fatalf("Secret not found at path: %s", fullPath)
	}

	// Print the secret
	fmt.Printf("Secret data: %v\n", secret.Data["data"])
}

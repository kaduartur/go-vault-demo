package vault

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/vault/api"
)

const vaultPath = "demo/credit-card/%s"
const vaultEncrypt = "demo/transit/encrypt/%s"
const vaultDecrypt = "demo/transit/decrypt/%s"

type Manager struct {
	client *api.Client
}

func NewManager(c *api.Client) *Manager {
	return &Manager{client: c}
}

func (vm *Manager) Write(key string, data map[string]interface{}) error {
	vm.setToken()
	if _, err := vm.client.Logical().Write(pathResolver(key), data); err != nil {
		log.Println("Vault write error", err)
		return err
	}
	return nil
}

func (vm *Manager) Read(key string) (map[string]interface{}, error) {
	vm.setToken()
	res, err := vm.client.Logical().Read(pathResolver(key))

	if err != nil {
		log.Println("Vault read error", err)
		return nil, err
	}
	if res == nil {
		return nil, nil
	}

	return res.Data, nil
}

func (vm *Manager) Delete(key string) error {
	vm.setToken()
	_, err := vm.client.Logical().Delete(pathResolver(key))
	if err != nil {
		log.Println("Vault delete error", err)
		return err
	}

	return nil
}

func (vm *Manager) Encrypt(data string) (string, error) {
	vm.setToken()
	path := fmt.Sprintf(vaultEncrypt, "credit_card")
	body := make(map[string]interface{})
	body["plaintext"] = base64.StdEncoding.EncodeToString([]byte(data))
	res, err := vm.client.Logical().Write(path, body)
	if err != nil {
		log.Println("Vault encrypt error", err)
		return "", err
	}
	return res.Data["ciphertext"].(string), nil
}

func (vm *Manager) setToken() {
	vm.client.SetToken(os.Getenv(api.EnvVaultToken))
}

func pathResolver(key string) string {
	return fmt.Sprintf(vaultPath, key)
}

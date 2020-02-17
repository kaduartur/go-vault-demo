package vault

import (
	"fmt"
	"log"
	"os"

	"github.com/hashicorp/vault/api"
)

const pathPattern = "%s/creds/%s"

var (
	vaultBackend = os.Getenv("VAULT_DB_BACKEND")
	vaultRole    = os.Getenv("VAULT_DB_ROLE")
)

type DatabaseManager struct {
	client *api.Client
}

func NewDatabaseManager(c *api.Client) *DatabaseManager {
	return &DatabaseManager{client: c}
}

func (d *DatabaseManager) CreateDynamicCredential() *api.Secret {
	d.setToken()
	res, err := d.client.Logical().Read(fmt.Sprintf(pathPattern, vaultBackend, vaultRole))
	if err != nil {
		log.Panicln(err)
	}

	_ = os.Setenv("DATABASE_USER", res.Data["username"].(string))
	_ = os.Setenv("DATABASE_PASSWORD", res.Data["password"].(string))

	return res
}

func (d *DatabaseManager) setToken() {
	d.client.SetToken(os.Getenv(api.EnvVaultToken))
}

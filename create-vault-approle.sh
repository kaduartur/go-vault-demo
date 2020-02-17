#!/bin/sh

C_DIR=$1
VAULT_HOST=$2

mkdir -p /tmp/vault

echo "Exporting vault vars..."
export VAULT_TOKEN="0f55f2f0-4e78-4331-8597-5b4f0abbcd08"
export VAULT_ADDR=$VAULT_HOST

echo "Installing vault cli 1.3.0..."
rm -rf /tmp/vault/vault
unzip "$C_DIR"/resources/vault_1.3.0_"$(uname -s)"_amd64.zip -d /tmp/vault/

# Enable secrets for paths
/tmp/vault/vault secrets enable -path=demo/credit-card generic
/tmp/vault/vault secrets enable -path=demo/transit transit
/tmp/vault/vault secrets enable database

# Create a vault policy
/tmp/vault/vault policy write demo_policy "$C_DIR"/resources/demo_policy.hcl

# Enable approle auth type
/tmp/vault/vault auth enable approle

# Add policy for authentication path
/tmp/vault/vault write auth/approle/role/demo_role policies=demo_policy period=15s

# Create a key for transit (encrypt and decrypt)
/tmp/vault/vault write -f demo/transit/keys/credit_card

# Configure PostgreSQL secrets engine
/tmp/vault/vault write database/config/demo_role \
  plugin_name=postgresql-database-plugin \
  allowed_roles=demo_role \
  connection_url="postgresql://{{username}}:{{password}}@postgresql:5432/vault_demo?sslmode=disable" \
  username="root" \
  password="root"

# Create a role
/tmp/vault/vault write database/roles/demo_role \
  db_name=demo_role \
  creation_statements="CREATE ROLE \"{{name}}\" WITH LOGIN PASSWORD '{{password}}' VALID UNTIL '{{expiration}}'; \
        GRANT SELECT, INSERT, UPDATE, DELETE ON ALL TABLES IN SCHEMA public TO \"{{name}}\";" \
  default_ttl=3h max_ttl=24h

rm -f /tmp/vault/role-id.txt
rm -f /tmp/vault/secret-id.txt

role_response=$(/tmp/vault/vault read -format=json auth/approle/role/demo_role/role-id)
echo "role_response $role_response"

role_id=$(echo "$role_response" | "$C_DIR"/jq-"$(uname -s)" -j '.data.role_id')
echo "role_id: $role_id"
eval echo "$role_id" >>/tmp/vault/role-id.txt

secret_response=$(/tmp/vault/vault write -force -format=json auth/approle/role/demo_role/secret-id)
echo "secret_response: $secret_response"

secret_id=$(echo "$secret_response" | "$C_DIR"/jq-"$(uname -s)" -j '.data.secret_id')
echo "secret_id: $secret_id"
eval echo "$secret_id" >>/tmp/vault/secret-id.txt

unset VAULT_TOKEN

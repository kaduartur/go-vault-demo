#!/bin/bash

./create-vault-approle.sh . http://0.0.0.0:8200

export VAULT_ADDR=http://localhost:8200
export VAULT_AUTHENTICATION=APPROLE
export VAULT_ROLE_ID=$(cat /tmp/vault/role-id.txt)
export VAULT_SECRET_ID=$(cat /tmp/vault/secret-id.txt)

export DATABASE_HOST=localhost
export DATABASE_PORT=5432
export DATABASE_USER=root
export DATABASE_PASSWORD=root
export DATABASE_NAME=vault_demo

go run server/cmd/server/main.go